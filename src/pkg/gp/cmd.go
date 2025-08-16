package gp

import "iter"

type GamespyCommandInfo struct {
	Command    string
	SubCommand int
	Data       iter.Seq[GamespyKV]
}

func (cmdInfo GamespyCommandInfo) Find(key string) GamespyKV {
	return FindInternal(key, cmdInfo.Data)
}

func FindInternal(key string, kvs iter.Seq[GamespyKV]) GamespyKV {

	for kv := range kvs {
		if kv.key != key {
			continue
		}

		return kv
	}

	return GamespyKV{}
}
