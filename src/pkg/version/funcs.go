package version

import (
	"fmt"
	"slices"
)

func GetFullTag() string {
	return fmt.Sprintf("%s.%s.%s.%s",
		GetVersion(),
		GetConfig(),
		GetLab(),
		GetIdentifier(),
	)
}

func GetPartialTag() string {
	return fmt.Sprintf("%s.%s.%s",
		GetConfig(),
		GetLab(),
		GetIdentifier(),
	)
}

func GetService() string {
	return DoNotTouch_Build_Service
}

func GetVersion() string {
	return fmt.Sprintf("%s.%s.%s.%s",
		DoNotTouch_Build_Version_Major,
		DoNotTouch_Build_Version_Minor,
		DoNotTouch_Build_Version_HotFix,
		DoNotTouch_Build_Version_Commit,
	)
}

func GetConfig() string {
	return fmt.Sprintf("%s%s.%s",
		DoNotTouch_Build_Cfg_OSArch,
		DoNotTouch_Build_Cfg_BldMode,
		DoNotTouch_Build_Cfg_OSName,
	)
}

func GetLab() string {
	labIden := ""
	if !slices.Contains(whitelistedBuildAccounts, DoNotTouch_Build_Lab_Username) {
		labIden = "(" + DoNotTouch_Build_Lab_Host + "\\" + DoNotTouch_Build_Lab_Username + ")"
	}

	return fmt.Sprintf("%s%s",
		DoNotTouch_Build_Lab_Name, labIden,
	)
}

func GetIdentifier() string {
	return fmt.Sprintf("%s.%s",
		DoNotTouch_Build_Iden_Timestamp,
		DoNotTouch_Build_Iden_DevFingerprint,
	)
}
