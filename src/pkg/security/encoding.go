package security

import (
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
)

func (s *Security) initb64wii() (outerr error) {
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("security: encoding: base64: %v", r)
		}
	}()

	// https://github.com/Retro-Rewind-Team/wfc-server/blob/main/common/encoding.go#L17
	s.Encoding.b64wii = base64.NewEncoding("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789.-").WithPadding('*')
	return nil
}

func (s *Security) initb32wii() (outerr error) {
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("security: encoding: base32: %v", r)
		}
	}()

	// https://github.com/Retro-Rewind-Team/wfc-server/blob/main/common/encoding.go#L20
	s.Encoding.b32wii = base32.NewEncoding("0123456789abcdefghijklmnopqrstuv")
	return nil
}

func (s SecEncoding) EncodeB64_Wii(data []byte) (outdata []byte, outerr error) {
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("security: encoding: base64: %v", r)
			outdata = nil
		}
	}()

	enc := []byte{}
	s.b64wii.Encode(enc, data)
	return enc, nil
}

func (s SecEncoding) DecodeB64_Wii(data []byte) (outdata []byte, outerr error) {
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("security: encoding: base64: %v", r)
			outdata = nil
		}
	}()

	dec := []byte{}
	n, err := s.b64wii.Decode(dec, data)
	if n == 0 {
		return nil, errors.New("security: encoding: base64: length was 0")
	}

	if err != nil {
		return nil, fmt.Errorf("security: encoding: base64: %w", err)
	}

	return dec, nil
}

func (s SecEncoding) EncodeB32_Wii(data []byte) (outdata []byte, outerr error) {
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("security: encoding: base32: %v", r)
			outdata = nil
		}
	}()

	enc := []byte{}
	s.b32wii.Encode(enc, data)
	return enc, nil
}

func (s SecEncoding) DecodeB32_Wii(data []byte) (outdata []byte, outerr error) {
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("security: encoding: base32: %v", r)
			outdata = nil
		}
	}()

	dec := []byte{}
	n, err := s.b32wii.Decode(dec, data)
	if n == 0 {
		return nil, errors.New("security: encoding: base32: length was 0")
	}

	if err != nil {
		return nil, fmt.Errorf("security: encoding: base32:: %w", err)
	}

	return dec, nil
}
