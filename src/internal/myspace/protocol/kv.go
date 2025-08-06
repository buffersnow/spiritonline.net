package protocol

import (
	"fmt"
	"iter"
	"slices"
	"strings"
)

type MySpaceKV struct {
	key   string
	value MySpaceData
}

func (kv MySpaceKV) Serialize() string {
	return fmt.Sprintf("\\%s\\%s", kv.key, kv.value)
}

func (kv MySpaceKV) Length() int {
	return len(kv.value.data)
}

func (kv MySpaceKV) Key() string {
	return kv.key
}

func (kv MySpaceKV) Value() MySpaceData {
	return kv.value
}

func SerializeKVDict(kvI iter.Seq[MySpaceKV]) string {
	final := ""
	for kv := range kvI {
		final += fmt.Sprintf("%s=%s\x1c", kv.key, kv.value)
	}
	return final
}

func DeserializeKVBatch(data string) iter.Seq[MySpaceKV] {

	splits := slices.DeleteFunc(strings.Split(data, "\\"), func(e string) bool {
		return e == "" // delete empty strings
	})

	return (func(yield func(MySpaceKV) bool) {

		for idx := 0; idx < len(splits); idx += 2 {

			val := ""
			if idx+1 < len(splits) {
				val = splits[idx+1]
			}

			if !yield(MySpaceKV{
				key:   splits[idx],
				value: MySpaceData{data: val},
			}) {
				return
			}
		}
	})
}

func DeserializeKVDict(data string) iter.Seq[MySpaceKV] {
	pairSplits := strings.SplitSeq(data, "\x1c")

	return func(yield func(MySpaceKV) bool) {
		for pair := range pairSplits {
			splits := strings.Split(pair, "=")

			for idx := 0; idx < len(splits); idx += 2 {

				value := ""
				if idx+1 < len(splits) {
					value = splits[idx+1]
				}

				if !yield(MySpaceKV{
					key:   splits[idx],
					value: MySpaceData{data: value},
				}) {
					return
				}
			}
		}
	}
}

func findInternal(key string, kvs iter.Seq[MySpaceKV]) MySpaceKV {

	for kv := range kvs {
		if kv.key != key {
			continue
		}

		return kv
	}

	return MySpaceKV{}
}

func (cbInfo MySpaceCallbackInfo) Find(key string) MySpaceKV {
	return findInternal(key, cbInfo.Body)
}

func (cmdInfo MySpaceCommandInfo) Find(key string) MySpaceKV {
	return findInternal(key, cmdInfo.Data)
}
