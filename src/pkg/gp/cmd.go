package gp

type GameSpyCommandInfo struct {
	Command    string
	SubCommand int
	Data       []GameSpyKV
}

func (cmdInfo GameSpyCommandInfo) Find(key string) GameSpyKV {
	return FindFromBundle(key, cmdInfo.Data)
}

func FindFromBundle(key string, kvs []GameSpyKV) GameSpyKV {

	for _, kv := range kvs {
		if kv.key != key {
			continue
		}

		return kv
	}

	return GameSpyKV{}
}
