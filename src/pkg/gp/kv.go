package gp

import (
	"fmt"
	"iter"
	"slices"
	"strings"
)

type GamespyKV struct {
	key   string
	value GamespyData
}

func (kv GamespyKV) Serialize() string {
	return fmt.Sprintf("\\%s\\%s", kv.key, kv.value)
}

func (kv GamespyKV) Length() int {
	return len(kv.value.data)
}

func (kv GamespyKV) Key() string {
	return kv.key
}

func (kv GamespyKV) Value() GamespyData {
	return kv.value
}

func SerializeKVDict(kvI iter.Seq[GamespyKV]) string {
	final := ""
	for kv := range kvI {
		final += fmt.Sprintf("%s=%s\x1c", kv.key, kv.value)
	}
	return final
}

func DeserializeKVBatch(data string) iter.Seq[GamespyKV] {

	splits := slices.DeleteFunc(strings.Split(data, "\\"), func(e string) bool {
		return e == "" // delete empty strings
	})

	return (func(yield func(GamespyKV) bool) {

		for idx := 0; idx < len(splits); idx += 2 {

			val := ""
			if idx+1 < len(splits) {
				val = splits[idx+1]
			}

			if !yield(GamespyKV{
				key:   splits[idx],
				value: GamespyData{data: val},
			}) {
				return
			}
		}
	})
}

func DeserializeKVDict(data string) iter.Seq[GamespyKV] {
	pairSplits := strings.SplitSeq(data, "\x1c")

	return func(yield func(GamespyKV) bool) {
		for pair := range pairSplits {
			splits := strings.Split(pair, "=")

			for idx := 0; idx < len(splits); idx += 2 {

				value := ""
				if idx+1 < len(splits) {
					value = splits[idx+1]
				}

				if !yield(GamespyKV{
					key:   splits[idx],
					value: GamespyData{data: value},
				}) {
					return
				}
			}
		}
	}
}
