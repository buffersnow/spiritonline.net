package security

import (
	"crypto/rc4"
	"crypto/sha1"
	"crypto/sha256"
	"hash/fnv"
)

type secCipher struct{}

var Cipher secCipher

func (secCipher) SHA1(data []byte) ([]byte, error) {
	hash := sha1.New()
	_, err := hash.Write(data)
	return hash.Sum(nil), err
}

func (secCipher) SHA2(data []byte) ([]byte, error) {
	hash := sha256.New()
	_, err := hash.Write(data)
	return hash.Sum(nil), err
}

func (secCipher) FNV(data string) (uint32, error) {
	hash := fnv.New32a()
	_, err := hash.Write([]byte(data))
	return hash.Sum32(), err
}

func (secCipher) RC4(pwd []byte, data []byte) ([]byte, error) {
	c, err := rc4.NewCipher(pwd)
	if err != nil {
		return nil, err
	}
	crypted := make([]byte, len(data))
	c.XORKeyStream(crypted, data)
	return crypted, nil
}
