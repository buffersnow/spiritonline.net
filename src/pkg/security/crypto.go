package security

import (
	"crypto/rc4"
	"fmt"
)

func (SecEncryption) RC4(pwd []byte, data []byte) (outdata []byte, outerr error) {
	c, err := rc4.NewCipher(pwd)
	if err != nil {
		return nil, fmt.Errorf("security: crypto: %w", err)
	}

	crypted := make([]byte, len(data))

	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("security: crypto: %v", r)
			outdata = nil
		}
	}()

	c.XORKeyStream(crypted, data) //~ why the fuck does this panic? terrible design
	return crypted, nil
}
