package gp

import (
	"fmt"
)

type GameSpyKV struct {
	key   string
	value GamespyData
}

func (kv GameSpyKV) Serialize() string {
	return fmt.Sprintf("\\%s\\%s", kv.key, kv.value)
}

func (kv GameSpyKV) Length() int {
	return len(kv.value.data)
}

func (kv GameSpyKV) Key() string {
	return kv.key
}

func (kv GameSpyKV) Value() GamespyData {
	return kv.value
}
