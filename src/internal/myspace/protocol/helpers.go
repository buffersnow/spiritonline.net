package protocol

import "strings"

func IdentifyProtocolRevision(ver int) string {

	if ver < 253 {
		return "v0 (Alpha)"
	} else if ver == 253 {
		return "v1"
	} else if ver > 253 && ver < 366 {
		return "v2"
	} else if ver >= 366 && ver < 404 {
		return "v3"
	} else if ver >= 404 && ver < 594 {
		return "v4"
	} else if ver >= 594 && ver < 673 {
		return "v5"
	} else if ver >= 673 && ver < 697 {
		return "v6"
	} else if ver >= 697 && ver < 812 {
		return "v7"
	}

	return "**Unknown**"
}

func EscapeString(data string) string {
	res := strings.ReplaceAll(data, "/", "/1")
	res = strings.ReplaceAll(res, "\\", "\\2")
	return res
}

func UnEscapeString(data string) string {
	res := strings.ReplaceAll(data, "/1", "/")
	res = strings.ReplaceAll(res, "\\2", "\\")
	return res
}
