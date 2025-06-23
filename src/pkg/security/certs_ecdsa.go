package security

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

type secCertsECDSA struct{}

func (secCertsECDSA) ReadPublicKey(pubKeyLocation string) (*ecdsa.PublicKey, error) {
	data, err := os.ReadFile(pubKeyLocation)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("no pem data found")
	} else if block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("expected \"PUBLIC KEY\" but got \"%s\"", block.Type)
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not an ECDSA public key")
	}

	return pKey, nil
}

func (secCertsECDSA) ReadPrivateKey(privKeyLocation string) (*ecdsa.PrivateKey, error) {
	data, err := os.ReadFile(privKeyLocation)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("no pem data found")
	} else if block.Type != "EC PRIVATE KEY" {
		return nil, fmt.Errorf("expected \"EC PRIVATE KEY\" but got \"%s\"", block.Type)
	}

	pKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return pKey, nil
}
