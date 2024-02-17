package auth

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"slices"
)

func Sign(message []byte, privateKey *ecdsa.PrivateKey) ([]byte, error) {
	hash := sha256.Sum256(message)

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	return slices.Concat[[]byte](hash[:], sig), nil
}

func Verify(message []byte, publicKey *ecdsa.PublicKey) bool {
	if len(message) <= 32 {
		return false
	}

	hash := message[:32]
	sig := message[32:]

	return ecdsa.VerifyASN1(publicKey, hash, sig)
}
