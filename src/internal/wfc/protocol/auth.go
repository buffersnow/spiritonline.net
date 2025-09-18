package protocol

import (
	"encoding/base64"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

/*migrated from json to cbor to make the token smaller*/
type AuthToken struct {
	WFCID      int64  `cbor:"wid"`
	GameCode   string `cbor:"gcd"`
	RegionID   byte   `cbor:"rid"`
	Serial     string `cbor:"csn"`
	FriendCode int64  `cbor:"cfc"`
	MAC        string `cbor:"mac"`
	IP         string `cbor:"ipa"`
	Challenge  string `cbor:"chg"`
}

func CreateToken(params AuthToken) (string, error) {
	token, err := cbor.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("cbor: %w", err)
	}

	return "NDS" + base64.StdEncoding.EncodeToString(token), nil
}
