package gp

import "iter"

type GamespyCommandInfo struct {
	Command    string
	SubCommand int
	Data       iter.Seq[GameSpyKV]
}

func (cmdInfo GamespyCommandInfo) Find(key string) GameSpyKV {
	return FindFromBundle(key, cmdInfo.Data)
}

func FindFromBundle(key string, kvs iter.Seq[GameSpyKV]) GameSpyKV {

	for kv := range kvs {
		if kv.key != key {
			continue
		}

		return kv
	}

	return GameSpyKV{}
}
