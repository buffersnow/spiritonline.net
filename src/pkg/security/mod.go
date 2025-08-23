package security

import (
	"crypto/ecdsa"
	"encoding/base32"
	"encoding/base64"

	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/version"
)

type Security struct {
	ECDSA    SecECDSA
	Hashing  SecHashing
	Crypto   SecEncryption
	Encoding SecEncoding
}

type SecECDSA struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

type SecHashing struct{}
type SecEncryption struct{}

type SecEncoding struct {
	b64wii *base64.Encoding
	b32wii *base32.Encoding
}

func New(bld *version.BuildTag, opt *settings.Options) (*Security, error) {
	sec := &Security{}

	//certsPath := fmt.Sprintf("%s/%s", opt.Runtime.CertsFolder, bld.GetService())
	tasks := []func() error{
		//func() error { return sec.readPublicKey(certsPath) },
		//func() error { return sec.readPrivateKey(certsPath) },
		func() error { return sec.initb64wii() },
		func() error { return sec.initb32wii() },
	}

	return sec, util.Batch(tasks)
}
