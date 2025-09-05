package gp

import (
	"fmt"
	"iter"
	"slices"
	"strings"
)

func DepickleDictionary(kvI iter.Seq[GameSpyKV]) string {
	final := ""
	for kv := range kvI {
		final += fmt.Sprintf("%s=%s\x1c", kv.key, kv.value)
	}
	return final
}

func PickleMessage(data string) iter.Seq[GameSpyKV] {

	splits := slices.DeleteFunc(strings.Split(data, "\\"), func(e string) bool {
		return e == "" //& delete empty strings
	})

	return (func(yield func(GameSpyKV) bool) {

		for idx := 0; idx < len(splits); idx += 2 {

			val := ""
			if idx+1 < len(splits) {
				val = splits[idx+1]
			}

			if !yield(GameSpyKV{
				key:   splits[idx],
				value: GamespyData{data: val},
			}) {
				return
			}
		}
	})
}

func PickleDictionary(data string) iter.Seq[GameSpyKV] {
	pairSplits := strings.SplitSeq(data, "\x1c")

	return func(yield func(GameSpyKV) bool) {
		for pair := range pairSplits {
			splits := strings.Split(pair, "=")

			for idx := 0; idx < len(splits); idx += 2 {

				value := ""
				if idx+1 < len(splits) {
					value = splits[idx+1]
				}

				if !yield(GameSpyKV{
					key:   splits[idx],
					value: GamespyData{data: value},
				}) {
					return
				}
			}
		}
	}
}
