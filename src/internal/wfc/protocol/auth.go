package protocol

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
)

type AuthToken struct {
	WFCID     int64  `json:"wfc_id"`
	GameCode  string `json:"game_cd"`
	RegionID  byte   `json:"region_id"`
	ConsoleID string `json:"console_id"`
	NandID    int64  `json:"nand_id"`
	MAC       string `json:"mac_addr"`
	IP        net.IP `json:"ip_addr"`
	Challenge string `json:"challenge"`
}

func CreateToken(params AuthToken) (string, error) {
	token, err := json.Marshal(params)
	if err != nil {
		return "", fmt.Errorf("encoding/json: %w", err)
	}

	return ("NDS" + base64.StdEncoding.EncodeToString(token)), nil
}
