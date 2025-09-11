package version

import "fmt"

//& This file only exists for ldflags version settings in run-service.sh
//! Variable available to LDFlags during compliation. DO NOT TOUCH!
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

//% Whitelisted Build Acounts
//% witsbla - Windows Test-Slave Build Account
//% lxmcbld - Linux Master Central Build Account

type buildTagVersion struct {
	major  string
	minor  string
	hotfix string
	commit string
}

type buildTagCfg struct {
	os_arch    string
	os_name    string
	build_mode string
}

type buildTagLab struct {
	name     string
	host     string
	username string
}

type buildTagIden struct {
	timestamp       string
	dev_fingerprint string
}

type BuildTag struct {
	service string
	version buildTagVersion
	cfg     buildTagCfg
	lab     buildTagLab
	iden    buildTagIden
}

func New() (*BuildTag, error) {
	bi := &BuildTag{
		service: DoNotTouch_Build_Service,
		version: buildTagVersion{
			major:  DoNotTouch_Build_Version_Major,
			minor:  DoNotTouch_Build_Version_Minor,
			hotfix: DoNotTouch_Build_Version_HotFix,
			commit: DoNotTouch_Build_Version_Commit,
		},
		cfg: buildTagCfg{
			os_arch:    DoNotTouch_Build_Cfg_OSArch,
			os_name:    DoNotTouch_Build_Cfg_OSName,
			build_mode: DoNotTouch_Build_Cfg_BldMode,
		},
		lab: buildTagLab{
			name:     DoNotTouch_Build_Lab_Name,
			host:     DoNotTouch_Build_Lab_Host,
			username: DoNotTouch_Build_Lab_Username,
		},
		iden: buildTagIden{
			timestamp:       DoNotTouch_Build_Iden_Timestamp,
			dev_fingerprint: DoNotTouch_Build_Iden_DevFingerprint,
		},
	}

	buffer := "- Welcome to spiritonline.net -> bFXServer/" + bi.GetService() + " v" + bi.GetVersion() + "\n"
	buffer += "Build Tag: " + bi.GetPartialTag()
	fmt.Printf("\033[38;5;61m%s\033[0m\n", buffer)

	return bi, nil
}
