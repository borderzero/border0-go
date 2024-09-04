package nacl

import (
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/curve25519"
	"golang.org/x/crypto/nacl/box"
)

const (
	// KeyLength is the length of keys (in bytes).
	KeyLength = 32

	// NonceLength is the the length of nonces (in bytes).
	NonceLength = 24
)

const (
	challengeLength = 32
)

// PublicKey represents a public key.
type PublicKey *[KeyLength]byte

// PrivateKey represents a private key.
type PrivateKey *[KeyLength]byte

// Nonce represents a nonce.
type Nonce *[NonceLength]byte

// Challenge represents something that can be solved and solutions can be checked against.
type Challenge interface {
	ToSolve() ([]byte, Nonce)
	IsSolution([]byte, Nonce) bool
}

// Service represents an entity capable of issuing challenges for others to solve.
type Service interface {
	PublicKey() PublicKey
	PrivateKey() PrivateKey
	NewChallengeForPeer(PublicKey) (Challenge, error)
}

// service is the default implementation of the Service interface.
type service struct {
	pub  PublicKey
	priv PrivateKey
}

// challenge is the default implementation of the Challenge interface.
type challenge struct {
	service Service

	peerPub PublicKey

	plaintext string
	nonce     Nonce
	sealed    []byte
}

// New returns a newly initialized, default Service implementation.
func New() (Service, error) {
	// NOTE: we could persist this in a secrets manager but with our current use
	// cases we don't have to. If we wanted clients to authenticate the server we
	// would pre-distribute the public key (but we have TLS for that). We only use
	// this package to authenticate clients, so we just generate a new key every time.
	var privateKey [KeyLength]byte
	if _, err := rand.Read(privateKey[:]); err != nil {
		return nil, fmt.Errorf("failed to generate private key for NaCl service: %v", err)
	}
	var publicKey [KeyLength]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)
	return &service{
		pub:  PublicKey(&publicKey),
		priv: PrivateKey(&privateKey),
	}, nil
}

// PublicKey returns the service's public key.
func (s *service) PublicKey() PublicKey { return s.pub }

// PrivateKey returns the service's private key.
func (s *service) PrivateKey() PrivateKey { return s.priv }

// NewChallengeForPeer creates a new challenge for a remote peer represented by
// the given public key. If the remote peer holds the corresponding private key,
// they will be able to solve the challenge.
func (s *service) NewChallengeForPeer(peerPub PublicKey) (Challenge, error) {
	challengeBytes, err := generateRandomChallenge(challengeLength)
	if err != nil {
		return nil, fmt.Errorf("failed to generate challenge message: %v", err)
	}
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}
	return &challenge{
		service: s,

		peerPub: peerPub,

		plaintext: string(challengeBytes),
		nonce:     nonce,
		sealed:    box.Seal(nil, challengeBytes, nonce, peerPub, s.priv),
	}, nil
}

// ToSolve returns the solvable challenge bytes (e.g. to be decrypted with their
// private key and re-encrypted with the server's public key).
func (c *challenge) ToSolve() ([]byte, Nonce) { return c.sealed, c.nonce }

// IsSolution checks if a given encrypted blob is a solution for the challenge.
func (c *challenge) IsSolution(soln []byte, nonce Nonce) bool {
	decrypted, ok := box.Open(nil, soln, nonce, c.peerPub, c.service.PrivateKey())
	if !ok {
		return false
	}
	if string(decrypted) == c.plaintext {
		return true
	}
	return false
}

// generateRandomChallenge generates a random challenge of a given length.
func generateRandomChallenge(n int) ([]byte, error) {
	challenge := make([]byte, n)
	if _, err := rand.Read(challenge); err != nil {
		return nil, fmt.Errorf("failed to generate challenge: %v", err)
	}
	return challenge, nil
}

// GenerateNonce generates a fixed number of random bytes.
func GenerateNonce() (Nonce, error) {
	var nonce [24]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}
	return &nonce, nil
}

// SolveChallenge attempts to solve a challenge with a given private key.
func SolveChallenge(challenge []byte, nonce Nonce, selfPriv PrivateKey, peerPub PublicKey) ([]byte, Nonce, error) {
	plaintextChallenge, ok := box.Open(nil, challenge, nonce, peerPub, selfPriv)
	if !ok {
		return nil, nil, errors.New("failed to decrypt challenge data with private key")
	}
	solnNonce, err := GenerateNonce()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce for challenge solution: %v", err)
	}
	return box.Seal(nil, plaintextChallenge, solnNonce, peerPub, selfPriv), solnNonce, nil
}
