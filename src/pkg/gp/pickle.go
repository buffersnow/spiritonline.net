package gp

import (
	"fmt"
	"slices"
	"strings"
)

func DepickleDictionary(kvs []GameSpyKV) string {
	final := ""
	for _, kv := range kvs {
		final += fmt.Sprintf("%s=%s\x1c", kv.key, kv.value)
	}
	return final
}

func PickleMessage(data string) []GameSpyKV {
	splits := slices.DeleteFunc(strings.Split(data, "\\"), func(e string) bool {
		return e == "" // delete empty strings
	})

	result := make([]GameSpyKV, 0, len(splits)/2)

	for idx := 0; idx < len(splits); idx += 2 {
		val := ""
		if idx+1 < len(splits) {
			val = splits[idx+1]
		}

		result = append(result, GameSpyKV{
			key:   splits[idx],
			value: GameSpyData{data: val},
		})
	}

	return result
}

func PickleDictionary(data string) []GameSpyKV {
	pairSplits := strings.Split(data, "\x1c")
	result := make([]GameSpyKV, 0, len(pairSplits))

	for _, pair := range pairSplits {
		splits := strings.Split(pair, "=")

		for idx := 0; idx < len(splits); idx += 2 {
			value := ""
			if idx+1 < len(splits) {
				value = splits[idx+1]
			}

			result = append(result, GameSpyKV{
				key:   splits[idx],
				value: GameSpyData{data: value},
			})
		}
	}

	return result
}
