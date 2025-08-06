package security

import (
	"crypto/ecdsa"
	"fmt"

	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/version"
)

type Security struct {
	ECDSA   SecECDSA
	Hashing SecHashing
	Crypto  SecEncryption
}

type SecECDSA struct {
	PublicKey  *ecdsa.PublicKey
	PrivateKey *ecdsa.PrivateKey
}

type SecHashing struct{}
type SecEncryption struct{}

func New(bld *version.BuildTag, opt *settings.Options) (*Security, error) {
	sec := &Security{}

	certsPath := fmt.Sprintf("%s/%s", opt.Runtime.CertsFolder, bld.GetService())
	tasks := []func() error{
		func() error { return sec.readPublicKey(certsPath) },
		func() error { return sec.readPrivateKey(certsPath) },
	}

	return sec, util.Batch(tasks)
}
