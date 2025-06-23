package version

/*
	This file only exists for ldflags version settings in run-service.sh
*/

// Variable available to LDFlags during compliation. DO NOT TOUCH!
var (
	DoNotTouch_Build_Service             string
	DoNotTouch_Build_Version_Major       string
	DoNotTouch_Build_Version_Minor       string
	DoNotTouch_Build_Version_HotFix      string
	DoNotTouch_Build_Version_Commit      string
	DoNotTouch_Build_Cfg_OSArch          string
	DoNotTouch_Build_Cfg_BldMode         string
	DoNotTouch_Build_Cfg_OSName          string
	DoNotTouch_Build_Lab_Name            string
	DoNotTouch_Build_Lab_Host            string
	DoNotTouch_Build_Lab_Username        string
	DoNotTouch_Build_Iden_Timestamp      string
	DoNotTouch_Build_Iden_DevFingerprint string
)

var whitelistedBuildAccounts = []string{
	"witsbla", /*Windows Test-Slave Build Account*/
	"lxmcbld", /*Linux Master Central Build Account*/
}

type BuildTag struct {
	Program_Name        string
	Program_Service     string
	Version_Major       string
	Version_Minor       string
	Version_HotFix      string
	Version_Commit      string
	Cfg_OSArch          string
	Cfg_BldMode         string
	Cfg_OSName          string
	Lab_Name            string
	Lab_Host            string
	Lab_Username        string
	Iden_Timestamp      string
	Iden_DevFingerprint string
}
