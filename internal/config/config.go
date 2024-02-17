package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

type config struct {
	Server server
	Key    key
}

var cfg config

type key struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  *ecdsa.PublicKey
}

type server struct {
	Addr string
}

func GetConfig() config {
	return cfg
}

func init() {
	pubkey := os.Getenv("public_key")
	if pubkey == "" {
		log.Println("error getting env public_key")
	} else {
		pubblock, _ := pem.Decode([]byte(pubkey))
		pubIKey, err := x509.ParsePKIXPublicKey(pubblock.Bytes)
		if err != nil {
			log.Panic("error loading public key", err)
		}

		publicKey, ok := pubIKey.(*ecdsa.PublicKey)
		if !ok {
			log.Panic("error casting public key to ecdsa.PublicKey")
		}

		cfg.Key.PublicKey = publicKey
	}

	privkey := os.Getenv("private_key")
	if privkey == "" {
		log.Println("error getting env private_key")
	} else {
		block, _ := pem.Decode([]byte(privkey))
		privateKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			log.Panic("error loading private key")
		}

		cfg.Key.PrivateKey = privateKey
	}

	addr := os.Getenv("goshare_addr")
	if addr != "" {
		log.Println("error getting env goshare_addr")
	} else {
		cfg.Server.Addr = addr
	}
}
