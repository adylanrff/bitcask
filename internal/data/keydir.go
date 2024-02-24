package data

import (
	"github.com/adylanrff/bitcask/pkg/types"
)

type Keydir map[types.Key]KeydirEntry

func (k *Keydir) Put(key types.Key, value KeydirEntry) {
	(*k)[key] = value
}

func (k *Keydir) Get(key types.Key) (KeydirEntry, bool) {
	val, ok := (*k)[key]
	return val, ok
}

func (k *Keydir) Delete(key types.Key) {
	delete((*k), key)
}

type KeydirEntry struct {
	FileID      string
	ValueSize   uint64
	ValueOffset uint64
	Timestamp   uint64
}
