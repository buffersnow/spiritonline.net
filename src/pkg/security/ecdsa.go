package security

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

func (s *Security) readPublicKey(location string) error {
	data, err := os.ReadFile(fmt.Sprintf("%s_public_key.pem", location))
	if err != nil {
		return fmt.Errorf("security: ecdsa: os: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("security: ecdsa: no pem data found")
	} else if block.Type != "PUBLIC KEY" {
		return fmt.Errorf("security: ecdsa: expected \"PUBLIC KEY\" but got \"%s\"", block.Type)
	}

	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("security: ecdsa: %w", err)
	}

	pKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("security: ecdsa: not an ECDSA public key")
	}
	s.ECDSA.PublicKey = pKey

	return nil
}

func (s *Security) readPrivateKey(location string) error {
	data, err := os.ReadFile(fmt.Sprintf("%s_private_key.pem", location))
	if err != nil {
		return fmt.Errorf("security: ecdsa: os: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return errors.New("security: ecdsa: no pem data found")
	} else if block.Type != "EC PRIVATE KEY" {
		return fmt.Errorf("security: ecdsa: expected \"EC PRIVATE KEY\" but got \"%s\"", block.Type)
	}

	pKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return fmt.Errorf("security: ecdsa: %w", err)
	}
	s.ECDSA.PrivateKey = pKey

	return nil
}
