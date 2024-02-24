package data

import (
	"encoding/binary"
	"hash/crc32"

	"github.com/adylanrff/bitcask/pkg/types"
)

type Entry struct {
	Crc       uint32 // 32bit crc checksum
	Timestamp uint64
	Key       types.Key
	Value     types.Value
}

func NewEntry(timestamp uint64, key types.Key, value types.Value) *Entry {
	crcPayload := make([]byte, 0)
	binary.BigEndian.AppendUint64(crcPayload, timestamp)
	binary.BigEndian.AppendUint32(crcPayload, uint32(len(key)))
	binary.BigEndian.AppendUint32(crcPayload, uint32(len(value)))
	crcPayload = append(crcPayload, key...)
	crcPayload = append(crcPayload, value...)

	return &Entry{
		Crc:       crc32.ChecksumIEEE(crcPayload),
		Timestamp: timestamp,
		Key:       key,
		Value:     value,
	}
}
