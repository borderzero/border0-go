package listen

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/borderzero/border0-go/client"
	"github.com/cenkalti/backoff/v4"
	"golang.org/x/crypto/ssh"
)

// Listener is a net.Listener that connects to a Border0 tunnel server and forwards
// connections to the local machine.
type Listener struct {
	inner        net.Listener
	apiClient    client.Requester
	authToken    string
	socketName   string
	tunnelServer string
	errChan      chan error
	readyChan    chan bool
}

const defaultTunnelServer = "tunnel.border0.com:22"

// New creates a new Listener with the given options.
func New(options ...Option) *Listener {
	l := &Listener{
		authToken:    os.Getenv("BORDER0_AUTH_TOKEN"),
		socketName:   os.Getenv("BORDER0_SOCKET_NAME"),
		tunnelServer: os.Getenv("BORDER0_TUNNEL_SERVER"),
		errChan:      make(chan error),
		readyChan:    make(chan bool),
	}
	if l.tunnelServer == "" {
		l.tunnelServer = defaultTunnelServer
	}
	for _, option := range options {
		option(l)
	}
	if l.apiClient == nil {
		l.apiClient = client.New(client.WithAuthToken(l.authToken))
	}
	return l
}

// Start connects to Border0 tunnel server and starts listening for connections.
// It will block until the listener is ready to accept connections.
func (l *Listener) Start() error {
	go l.connectTunnel()

	select {
	case err := <-l.errChan:
		return err
	case <-l.readyChan:
	}

	go func() {
		for err := range l.errChan {
			fmt.Println("Border0 listener:", err)
		}
	}()

	return nil
}

// Accept waits for and returns the next connection to the listener.
func (l *Listener) Accept() (net.Conn, error) {
	conn, err := l.inner.Accept()
	if err != nil {
		if err == io.EOF {
			l.errChan <- fmt.Errorf("listener closed and reconnecting")
			<-l.readyChan
			return l.Accept()
		} else {
			return nil, err
		}
	}
	return conn, nil
}

// Close closes the listener.
func (l *Listener) Close() error {
	return l.inner.Close()
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.inner.Addr()
}

func (l *Listener) connectTunnel() {
	claims, err := l.apiClient.TokenClaims()
	if err != nil {
		l.errChan <- fmt.Errorf("failed to get api token claims: %v", err)
		return
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		l.errChan <- errors.New("can't find claim for user_id")
		return
	}

	if err := l.createSocket(); err != nil {
		l.errChan <- err
		return
	}

	keyPair, err := generateSSHKeyPair()
	if err != nil {
		l.errChan <- err
		return
	}

	retries := backoff.NewExponentialBackOff()
	err = backoff.Retry(func() error {
		signer, hostKey, err := l.sshCert(keyPair)
		if err != nil {
			return fmt.Errorf("failed to get ssh cert: %w", err)
		}

		sshConfig := &ssh.ClientConfig{
			User:            strings.ReplaceAll(userID, "-", ""),
			HostKeyCallback: ssh.FixedHostKey(hostKey),
			Timeout:         10 * time.Second,
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		}

		if err := l.connect(sshConfig, retries); err != nil {
			return err
		}

		return nil
	}, retries)

	if err != nil {
		l.errChan <- fmt.Errorf("error connecting to server: %v", err)
		log.Fatal(err)
		return
	}
}

func (l *Listener) connect(config *ssh.ClientConfig, retries backoff.BackOff) error {
	sshClient, err := ssh.Dial("tcp", l.tunnelServer, config)
	if err != nil {
		return err
	}
	defer sshClient.Close()

	done := make(chan bool, 1)
	defer func() {
		done <- true
	}()
	go l.keepAlive(sshClient, done)

	l.inner, err = sshClient.Listen("tcp", "localhost:0")
	if err != nil {
		return fmt.Errorf("failed to open listener on tunnel server: %w", err)
	}
	defer l.inner.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %v", err)
	}
	defer session.Close()

	session.Stdout = os.Stdout

	var modes ssh.TerminalModes
	if err := session.RequestPty("xterm-256color", 80, 40, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed %v", err)
	}

	if err := session.Shell(); err != nil {
		return err
	}

	// reset the backoff timer only when the session is finished or has error returned
	// by using defer, it ensures the reset gets called after session.Wait() and before
	// the return (with either an error or nil)
	defer retries.Reset()

	l.readyChan <- true
	if err := session.Wait(); err != nil {
		return err
	}

	return nil
}

func (l *Listener) keepAlive(sshClient *ssh.Client, done chan bool) {
	t := time.NewTicker(10 * time.Second)
	max := 4
	n := 0

	defer t.Stop()

	for {
		select {
		case <-done:
			return
		case <-t.C:
			aliveChan := make(chan bool, 1)

			go func() {
				_, _, err := sshClient.SendRequest("keepalive@openssh.com", true, nil)
				if err != nil {
					aliveChan <- false
				} else {
					aliveChan <- true
				}
			}()

			select {
			case <-time.After(5 * time.Second):
				n++
			case alive := <-aliveChan:
				if !alive {
					n++
				} else {
					n = 0
				}
			}

			if n >= max {
				log.Println("ssh keepalive timeout, disconnecting")
				sshClient.Close()
				return
			}
		}
	}
}

func (l *Listener) createSocket() error {
	if socket, err := l.apiClient.Socket(context.Background(), l.socketName); err != nil {
		if client.NotFound(err) {
			// socket doesn't exist, let's create it first
			socket = new(client.Socket)
			socket.Name = l.socketName
			socket.SocketType = "http"

			if _, err = l.apiClient.CreateSocket(context.Background(), socket); err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (l *Listener) sshCert(keyPair *sshKeyPair) (signer ssh.Signer, hostKey ssh.PublicKey, err error) {
	keyToSign := client.SocketKeyToSign{
		SSHPublicKey: strings.TrimRight(string(keyPair.publicKey), "\n"),
	}

	signedKey, err := l.apiClient.SignSocketKey(context.Background(), l.socketName, &keyToSign)
	if err != nil {
		return nil, nil, fmt.Errorf("error: %v", err)
	}

	if signedKey.SignedSSHCert == "" {
		return nil, nil, fmt.Errorf("error: Unable to get signed key from Server")
	}

	certData := []byte(signedKey.SignedSSHCert)
	pubcert, _, _, _, err := ssh.ParseAuthorizedKey(certData)
	if err != nil {
		return nil, nil, fmt.Errorf("error: %v", err)
	}

	cert, ok := pubcert.(*ssh.Certificate)
	if !ok {
		return nil, nil, fmt.Errorf("error failed to cast to certificate: %v", err)
	}

	clientKey, _ := ssh.ParsePrivateKey(keyPair.privateKey)
	certSigner, err := ssh.NewCertSigner(cert, clientKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create signer: %v", err)
	}

	hostKeyBytes, err := base64.StdEncoding.DecodeString(signedKey.HostKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode hostkey %v", err)
	}
	if hostKey, err = ssh.ParsePublicKey(hostKeyBytes); err != nil {
		return nil, nil, fmt.Errorf("failed to parse hostkey %v", err)
	}

	return certSigner, hostKey, nil
}

type sshKeyPair struct {
	privateKey []byte
	publicKey  []byte
}

func generateSSHKeyPair() (*sshKeyPair, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	var privBuf bytes.Buffer
	parsed, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, err
	}

	privPEM := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: parsed,
	}

	if err := pem.Encode(&privBuf, privPEM); err != nil {
		return nil, err
	}

	pub, err := ssh.NewPublicKey(&priv.PublicKey)
	if err != nil {
		return nil, err
	}

	return &sshKeyPair{
		privateKey: bytes.TrimSpace(privBuf.Bytes()),
		publicKey:  bytes.TrimSpace(ssh.MarshalAuthorizedKey(pub)),
	}, nil
}
