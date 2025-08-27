package protocol

import "time"

//$ https://github.com/WiiLink24/wfc-server/blob/main/nas/auth.go#L410

func GetDateTime() string {
	return time.Now().Format("20060102150405")
}
