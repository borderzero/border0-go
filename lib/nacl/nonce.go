package nacl

import (
	"crypto/rand"
	"fmt"
)

// NonceLength is the the length of nonces (in bytes).
const NonceLength = 24

// Nonce represents a nonce.
type Nonce *[NonceLength]byte

// GenerateNonce generates a fixed number of random bytes.
func GenerateNonce() (Nonce, error) {
	var nonce [NonceLength]byte
	if _, err := rand.Read(nonce[:]); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}
	return &nonce, nil
}
