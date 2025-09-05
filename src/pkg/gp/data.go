package gp

import (
	"encoding/base64"
	"iter"
	"strconv"
)

type GamespyData struct {
	data string
}

func (g GamespyData) DataBlock() []byte {
	return []byte(g.data)
}

func (g GamespyData) String() string {
	return g.data
}

func (g GamespyData) Integer() (int, error) {
	return strconv.Atoi(g.data)
}

func (g GamespyData) Dictionary() iter.Seq[GameSpyKV] {
	return PickleDictionary(g.data)
}

func (g GamespyData) Base64() string {
	return base64.StdEncoding.EncodeToString([]byte(g.data))
}
