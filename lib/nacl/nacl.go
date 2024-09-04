package nacl

import (
	"fmt"

	"golang.org/x/crypto/nacl/box"
)

// Service represents an entity capable of issuing challenges for others to solve.
type Service interface {
	PrivateKey() *PrivateKey
	NewChallengeForPeer(*PublicKey) (Challenge, error)
}

// service is the default implementation of the Service interface.
type service struct {
	privateKey *PrivateKey
}

// New returns a newly initialized, default Service implementation.
func New(privateKey *PrivateKey) (Service, error) { return &service{privateKey: privateKey}, nil }

// PrivateKey returns the service's private key.
func (s *service) PrivateKey() *PrivateKey { return s.privateKey }

// NewChallengeForPeer creates a new challenge for a remote peer represented by
// the given public key. If the remote peer holds the corresponding private key,
// they will be able to solve the challenge.
func (s *service) NewChallengeForPeer(peersPublicKey *PublicKey) (Challenge, error) {
	challengeBytes, err := generateRandomChallenge()
	if err != nil {
		return nil, fmt.Errorf("failed to generate challenge message: %v", err)
	}
	nonce, err := GenerateNonce()
	if err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %v", err)
	}
	return &challenge{
		service: s,

		peersPublicKey: peersPublicKey,

		plaintext: string(challengeBytes),
		nonce:     nonce,
		sealed:    box.Seal(nil, challengeBytes, nonce, peersPublicKey.Raw(), s.privateKey.Raw()),
	}, nil
}
