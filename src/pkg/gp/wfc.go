package gp

import (
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/security"
	"github.com/fxamacker/cbor/v2"
	"github.com/luxploit/red"
)

/*migrated from json to cbor to make the token smaller*/
type WFCAuthToken struct {
	WFCID      int64     `cbor:"wid"`
	GameCode   string    `cbor:"gcd"`
	RegionID   byte      `cbor:"rid"`
	Serial     string    `cbor:"csn"`
	FriendCode int64     `cbor:"cfc"`
	MAC        string    `cbor:"mac"`
	IP         string    `cbor:"ipa"`
	Challenge  string    `cbor:"chg"`
	IssueTime  time.Time `cbor:"ite"`
}

func PickleWFCToken(params WFCAuthToken) (string, error) {

	sec, err := red.Locate[security.Security]()
	if err != nil {
		return "", fmt.Errorf("gp: %w", err)
	}

	params.IssueTime = time.Now()
	token, err := cbor.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("gp: cbor: %w", err)
	}

	b64, err := sec.Encoding.EncodeB64_Std(token)
	if err != nil {
		return "", fmt.Errorf("gp: %w", err)
	}

	return "NDS" + string(b64), nil
}

func DepickleWFCToken(wfctoken []byte) (*WFCAuthToken, error) {

	sec, err := red.Locate[security.Security]()
	if err != nil {
		return nil, fmt.Errorf("gp: %w", err)
	}

	b64, err := sec.Encoding.DecodeB64_Std(wfctoken)
	if err != nil {
		return nil, fmt.Errorf("gp: %w", err)
	}

	var token *WFCAuthToken
	err = cbor.Unmarshal(b64, &token)
	if err != nil {
		return nil, fmt.Errorf("gp: cbor: %w", err)
	}

	return token, nil
}
