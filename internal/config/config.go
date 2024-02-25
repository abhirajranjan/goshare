package config

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"reflect"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/pkg/errors"
)

type config struct {
	Server server
	Key    key
	Pubsub pubsub
	Otp    otp
}

var cfg config
var fnMap = make(map[reflect.Type]env.ParserFunc)

type key struct {
	PrivateKey *ecdsa.PrivateKey `env:"key.private" envDefault:"privatekey"`
	PublicKey  *ecdsa.PublicKey  `env:"key.pub" envDefault:"pubkey"`
}

type server struct {
	Addr string `env:"srv.addr" envDefault:":8080"`
}

type pubsub struct {
	TTL time.Duration `env:"pubsub.ttl" envDefault:"1m"`
}

type otp struct {
	ValidDuration time.Duration `env:"otp.valid_duration" envDefault:"30s"`
}

func GetConfig() config {
	return cfg
}

func init() {
	fnMap[reflect.TypeFor[ecdsa.PrivateKey]()] = func(v string) (interface{}, error) {
		privateKeyByte, err := os.ReadFile(v)
		if err != nil {
			return nil, err
		}

		block, _ := pem.Decode([]byte(privateKeyByte))
		privateKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return nil, errors.Wrap(err, "init.private_key")
		}

		return *privateKey, nil
	}

	fnMap[reflect.TypeFor[ecdsa.PublicKey]()] = func(v string) (interface{}, error) {
		publicKeyByte, err := os.ReadFile(v)
		if err != nil {
			return nil, errors.Wrap(err, "init.public_key")
		}

		block, _ := pem.Decode([]byte(publicKeyByte))
		pubIKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return nil, errors.Wrap(err, "init.public_key")
		}

		publicKey, ok := pubIKey.(*ecdsa.PublicKey)
		if !ok {
			return nil, errors.Errorf("cannot cast publicKey[%T] to *ecdsa.PublicKey", publicKey)
		}

		return *publicKey, nil
	}

	fnMap[reflect.TypeFor[time.Duration]()] = func(v string) (interface{}, error) {
		return time.ParseDuration(v)
	}
}

func init() {
	if err := env.ParseWithOptions(&cfg, env.Options{
		Environment:     env.ToMap(os.Environ()),
		RequiredIfNoDef: true,
		FuncMap:         fnMap,
	}); err != nil {
		panic(err)
	}
}
