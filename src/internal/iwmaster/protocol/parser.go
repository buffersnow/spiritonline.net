package protocol

import (
	"slices"
	"strings"
)

func PickleMessage(data []byte) IWCommandInfo {
	wci := IWCommandInfo{}

	if !slices.Equal(data[0:4], []byte{0xff, 0xff, 0xff, 0xff}) {
		return wci
	}
	data = data[4 : len(data)-3] //@ TODO: this will break older games

	splits := strings.FieldsFunc(string(data), func(r rune) bool {
		return r == ' ' || r == '\n'
	})

	splits = slices.DeleteFunc(splits, func(e string) bool {
		return e == ""
	})

	if len(splits) == 0 || len(splits[0]) == 0 {
		return wci
	}

	wci.Command = splits[0]
	if len(splits) > 1 {
		wci.Data = splits[1:]
	}

	return wci
}
