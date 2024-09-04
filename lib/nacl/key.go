package nacl

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"golang.org/x/crypto/curve25519"
)

// KeyLength is the length of keys (in bytes).
const KeyLength = 32

// PublicKey represents a public key.
type PublicKey struct {
	raw *[KeyLength]byte
}

// PrivateKey represents a private key.
type PrivateKey struct {
	raw *[KeyLength]byte
	pub *PublicKey
}

// GenerateKey generates a new PrivateKey.
func GenerateKey() (*PrivateKey, error) {
	var privateKey [KeyLength]byte
	if _, err := rand.Read(privateKey[:]); err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	var publicKey [KeyLength]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)
	return &PrivateKey{
		raw: &privateKey,
		pub: &PublicKey{raw: &publicKey},
	}, nil
}

// ParsePrivateKey parses raw bytes onto a PrivateKey.
func ParsePrivateKey(raw []byte) (*PrivateKey, error) {
	if len(raw) != KeyLength {
		return nil, fmt.Errorf("invalid key length: expected %d but got %d bytes", KeyLength, len(raw))
	}
	privateKey := [KeyLength]byte(raw)
	var publicKey [KeyLength]byte
	curve25519.ScalarBaseMult(&publicKey, &privateKey)
	return &PrivateKey{
		raw: &privateKey,
		pub: &PublicKey{raw: &publicKey},
	}, nil
}

// ParsePrivateKeyB64 parses base64-encoded bytes onto a PrivateKey.
func ParsePrivateKeyB64(b64 string) (*PrivateKey, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to base64-decode given string: %v", err)
	}
	return ParsePrivateKey(raw)
}

// ParsePublicKey parses raw bytes onto a PublicKey.
func ParsePublicKey(raw []byte) (*PublicKey, error) {
	if len(raw) != KeyLength {
		return nil, fmt.Errorf("invalid key length: expected %d but got %d bytes", KeyLength, len(raw))
	}
	publicKey := [KeyLength]byte(raw)
	return &PublicKey{raw: &publicKey}, nil
}

// ParsePublicKeyB64 parses base64-encoded bytes onto a PublicKey.
func ParsePublicKeyB64(b64 string) (*PublicKey, error) {
	raw, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, fmt.Errorf("failed to base64-decode given string: %v", err)
	}
	return ParsePublicKey(raw)
}

// Raw returns the address of the raw public key bytes.
func (p *PublicKey) Raw() *[KeyLength]byte {
	return p.raw
}

// Public returns the public key corresponding to the private key.
func (p *PrivateKey) Public() *PublicKey {
	return p.pub
}

// Raw returns the address of the raw private key bytes.
func (p *PrivateKey) Raw() *[KeyLength]byte {
	return p.raw
}
