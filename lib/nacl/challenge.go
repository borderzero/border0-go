package nacl

import (
	"crypto/rand"
	"errors"
	"fmt"

	"golang.org/x/crypto/nacl/box"
)

// challengeLength is the length of generated challenges.
const challengeLength = 32

// Challenge represents something that can be solved and solutions can be checked against.
type Challenge interface {
	ToSolve() ([]byte, Nonce)
	IsSolution([]byte, Nonce) bool
}

// challenge is the default implementation of the Challenge interface.
type challenge struct {
	service Service

	peersPublicKey *PublicKey

	plaintext string
	nonce     Nonce
	sealed    []byte
}

// ToSolve returns the solvable challenge bytes (e.g. to be decrypted with their
// private key and re-encrypted with the server's public key).
func (c *challenge) ToSolve() ([]byte, Nonce) { return c.sealed, c.nonce }

// IsSolution checks if a given encrypted blob is a solution for the challenge.
func (c *challenge) IsSolution(soln []byte, nonce Nonce) bool {
	decrypted, ok := box.Open(nil, soln, nonce, c.peersPublicKey.Raw(), c.service.PrivateKey().Raw())
	if !ok {
		return false
	}
	if string(decrypted) == c.plaintext {
		return true
	}
	return false
}

// SolveChallenge attempts to solve a challenge with a given private key.
func SolveChallenge(challenge []byte, nonce Nonce, peersPublicKey *PrivateKey, privateKey *PublicKey) ([]byte, Nonce, error) {
	plaintextChallenge, ok := box.Open(nil, challenge, nonce, peersPublicKey.Raw(), privateKey.Raw())
	if !ok {
		return nil, nil, errors.New("failed to decrypt challenge data with private key")
	}
	solnNonce, err := GenerateNonce()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate nonce for challenge solution: %v", err)
	}
	return box.Seal(nil, plaintextChallenge, solnNonce, peersPublicKey.Raw(), privateKey.Raw()), solnNonce, nil
}

// generateRandomChallenge generates a random challenge of a given length.
func generateRandomChallenge() ([]byte, error) {
	challenge := make([]byte, challengeLength)
	if _, err := rand.Read(challenge); err != nil {
		return nil, fmt.Errorf("failed to generate challenge: %v", err)
	}
	return challenge, nil
}
