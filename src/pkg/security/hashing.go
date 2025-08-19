package security

import (
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"hash/fnv"
)

func (SecHashing) SHA1(data []byte) ([]byte, error) {
	hash := sha1.New()

	if _, err := hash.Write(data); err != nil {
		return nil, fmt.Errorf("security: hash: sha1: %w", err)
	}

	return hash.Sum(nil), nil
}

func (SecHashing) SHA2(data []byte) ([]byte, error) {
	hash := sha256.New()

	if _, err := hash.Write(data); err != nil {
		return nil, fmt.Errorf("security: hash: sha2: %w", err)
	}

	return hash.Sum(nil), nil
}

func (SecHashing) FNV(data []byte) (uint32, error) {
	hash := fnv.New32a()

	if _, err := hash.Write(data); err != nil {
		return 0, fmt.Errorf("security: hash: fnv: %w", err)
	}

	return hash.Sum32(), nil
}
