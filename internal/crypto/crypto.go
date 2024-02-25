package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
)

var baseEnc = base64.RawStdEncoding.WithPadding(base64.NoPadding)

type crypto struct {
	pubKey     *ecdsa.PublicKey
	privateKey *ecdsa.PrivateKey
}

type CryptoOpts func(*crypto)

func NewCrypto(pubKey *ecdsa.PublicKey, privateKey *ecdsa.PrivateKey) (c crypto) {
	c.privateKey = privateKey
	c.pubKey = pubKey
	return
}

func (c crypto) Sign(message []byte) ([]byte, error) {
	if len(message) >= 1<<8 {
		return nil, errors.New("too long message (> 1 byte)")
	}

	hash := sha256.Sum256(message)
	sig, err := ecdsa.SignASN1(rand.Reader, c.privateKey, hash[:])
	if err != nil {
		return nil, err
	}

	token := make([]byte, 0, 1+len(message)+len(hash)+len(sig))
	token = append(token, byte(len(message)))
	token = append(token, message...)
	token = append(token, hash[:]...)
	token = append(token, sig...)

	buf := make([]byte, baseEnc.EncodedLen(len(token)))
	baseEnc.Encode(buf, token)

	return buf, nil
}

func (c crypto) Verify(token []byte) bool {
	buf := decode(token)
	if len(buf) == 0 {
		return false
	}

	messageEnd := 1 + buf[0]
	hash := buf[messageEnd : 32+messageEnd]
	sig := buf[32+messageEnd:]

	return ecdsa.VerifyASN1(c.pubKey, hash, sig)
}

func (c crypto) GetMessage(token []byte) []byte {
	buf := decode(token)
	if len(buf) == 0 {
		return nil
	}

	messageEnd := 1 + buf[0]
	return buf[1 : 1+messageEnd]
}

func decode(token []byte) []byte {
	buf := make([]byte, baseEnc.DecodedLen(len(token)))
	if _, err := baseEnc.Decode(buf, token); err != nil {
		log.Printf("error decoding token: %s\n", err)
		return nil
	}

	if len(buf) <= 33 {
		return nil
	}

	return buf
}
