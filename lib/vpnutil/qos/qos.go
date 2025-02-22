package qos

import (
	"encoding/binary"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/borderzero/border0-go/lib/nacl"
	"github.com/borderzero/border0-go/lib/types/pointer"
	"golang.org/x/crypto/nacl/box"
)

const (
	magicCookie = 0x2112A443 // one more than the STUN magic cookie

	// sizes
	probeIdLen = 4
	sentAtLen  = 4

	// probe message offsets
	cookieStartOffset        = 0
	cookieEndOffset          = cookieStartOffset + 4
	msgTypeStartOffset       = cookieEndOffset
	msgTypeEndOffset         = msgTypeStartOffset + 1
	publicKeyStartOffset     = msgTypeEndOffset
	publicKeyEndOffset       = publicKeyStartOffset + nacl.KeyLength
	nonceStartOffset         = publicKeyEndOffset
	nonceEndOffset           = nonceStartOffset + nacl.NonceLength
	encryptedBodyStartOffset = nonceEndOffset
	encryptedBodyEndOffset   = encryptedBodyStartOffset + probeIdLen + sentAtLen + box.Overhead
	qosMessageSize           = encryptedBodyEndOffset

	// decrypted/plaintext probe body offsets
	probeIdStartBodyOffset = 0
	probeIdEndBodyOffset   = probeIdStartBodyOffset + probeIdLen
	sentAtStartBodyOffset  = probeIdEndBodyOffset
	sentAtEndBodyOffset    = sentAtStartBodyOffset + sentAtLen
)

const (
	// MessageTypeRequest is the message type for requests.
	MessageTypeRequest = byte(0x01)

	// MessageTypeResponse is the message type for responses.
	MessageTypeResponse = byte(0x02)
)

var bin = binary.BigEndian

// Message represents a QOS message, which may be a request or response.
type Message struct {
	key *nacl.PrivateKey

	cookie  uint32
	remote  *nacl.PublicKey
	msgtype byte
	probeId uint32
	sentAt  time.Time
}

// NewRequest builds a new QOS request for the given remote peer key.
func NewRequest(priv *nacl.PrivateKey, remote *nacl.PublicKey) *Message {
	return &Message{
		key: priv,

		cookie:  magicCookie,
		msgtype: MessageTypeRequest,
		remote:  remote,
		probeId: rand.Uint32(),
		sentAt:  time.Now(),
	}
}

// NewResponse builds a new QOS response for the given request.
func NewResponse(priv *nacl.PrivateKey, req *Message) *Message {
	return &Message{
		key: priv,

		cookie:  magicCookie,
		msgtype: MessageTypeResponse,
		remote:  req.remote,
		probeId: req.probeId,
		sentAt:  req.sentAt,
	}
}

// Encode encodes a Message onto wire-ready bytes.
func (m *Message) Encode() []byte {
	nonce, err := nacl.GenerateNonce()
	if err != nil {
		// this can only fail if there is an issue with the underlying RNG.
		panic(fmt.Errorf("failed to generate nonce for message encoding: %v", err))
	}

	var buf []byte
	buf = bin.AppendUint32(buf, magicCookie)
	buf = append(buf, m.msgtype)
	buf = append(buf, m.key.Public().Raw()[:]...)
	buf = append(buf, nonce[:]...)

	var body []byte
	body = bin.AppendUint32(body, m.probeId)
	body = bin.AppendUint32(body, uint32(m.sentAt.UnixMicro()))

	return box.Seal(buf, body, nonce, m.remote.Raw(), m.key.Raw())
}

// Remote returns the key of the remote peer of the message.
func (m *Message) Remote() *nacl.PublicKey { return m.remote }

// MessageType returns the type of the message.
func (m *Message) MessageType() byte { return m.msgtype }

// ProbeID returns the probe id of the message.
func (m *Message) ProbeID() uint32 { return m.probeId }

// ParseQOSMessage takes in bytes for a packet received over the network,
// it returns the parsed message, true if the bytes indeed corresponded to
// a message, and an error if and only if the packet was for a QOS message
// but it failed to be parsed.
func ParseQOSMessage(key *nacl.PrivateKey, pck []byte) (*Message, bool, error) {
	if len(pck) < qosMessageSize {
		return nil, false, nil
	}

	cookie := bin.Uint32(pck[cookieStartOffset:cookieEndOffset])
	if cookie != magicCookie {
		return nil, false, nil
	}

	msgType := pck[msgTypeStartOffset]

	pub, err := nacl.ParsePublicKey(pck[publicKeyStartOffset:publicKeyEndOffset])
	if err != nil {
		return nil, true, fmt.Errorf("failed to parse public key: %v", err)
	}

	body, ok := box.Open(
		nil,
		pck[encryptedBodyStartOffset:encryptedBodyEndOffset],
		pointer.To([nacl.NonceLength]byte(pck[nonceStartOffset:nonceEndOffset])),
		pub.Raw(),
		key.Raw(),
	)
	if !ok {
		return nil, true, fmt.Errorf("failed to authenticate QOS message: message not an NaCl box sealed message with the key %s", pub.String())
	}

	return &Message{
		remote:  pub,
		cookie:  cookie,
		msgtype: msgType,
		probeId: bin.Uint32(body[probeIdStartBodyOffset:probeIdEndBodyOffset]),
		sentAt:  time.UnixMicro(int64(bin.Uint32(body[sentAtStartBodyOffset:sentAtEndBodyOffset]))),
	}, true, nil
}
