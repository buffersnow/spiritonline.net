package protocol

import (
	"encoding/base64"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

/*migrated from json to cbor to make the token smaller*/
type AuthToken struct {
	WFCID     int64  `cbor:"wid"`
	GameCode  string `cbor:"gcd"`
	RegionID  byte   `cbor:"rid"`
	ConsoleID string `cbor:"cid"`
	MAC       string `cbor:"mac"`
	IP        string `cbor:"ipa"`
	Challenge string `cbor:"chg"`
}

func CreateToken(params AuthToken) (string, error) {
	token, err := cbor.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("cbor: %w", err)
	}

	fmt.Printf("%s\n", string(token))
	fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(token))

	return "NDS" + base64.StdEncoding.EncodeToString(token), nil
}
