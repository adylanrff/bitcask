package codec

import (
	"encoding/binary"
	"io"

	"github.com/adylanrff/bitcask/internal/data"
	"github.com/adylanrff/bitcask/pkg/types"
)

// Encoding:
// [crc(4),keySize(4),valueSize(8),actualKey(?),actualValue(?)]

const (
	keySizeByteSize   = 4 // 32bit
	valueSizeByteSize = 8 // 64bit
	crcSizeByteSize   = 4 // 32bit
)

type binaryCodec struct{}

// Marshal implements Codec.
func (b *binaryCodec) Encode(w io.Writer, entry *data.Entry) (uint64, error) {
	buf := make([]byte, crcSizeByteSize+keySizeByteSize+valueSizeByteSize)
	binary.BigEndian.PutUint32(buf[:crcSizeByteSize], entry.Crc)
	binary.BigEndian.PutUint32(buf[crcSizeByteSize:crcSizeByteSize+keySizeByteSize], uint32(len(entry.Key)))
	binary.BigEndian.PutUint64(buf[crcSizeByteSize+keySizeByteSize:crcSizeByteSize+keySizeByteSize+valueSizeByteSize], uint64(len(entry.Value)))

	n, err := w.Write(buf)
	if err != nil {
		return 0, err
	}

	totalSize := int64(n)

	n, err = w.Write([]byte(entry.Key))
	if err != nil {
		return 0, err
	}
	totalSize += int64(n)

	n, err = w.Write(entry.Value)
	if err != nil {
		return 0, err
	}
	totalSize += int64(n)

	return uint64(totalSize), nil
}

// Unmarshal implements Codec.
func (*binaryCodec) Decode(r io.Reader) (*data.Entry, uint64, error) {
	buf := make([]byte, crcSizeByteSize+keySizeByteSize+valueSizeByteSize)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, 0, err
	}

	crc := binary.BigEndian.Uint32(buf[:crcSizeByteSize])
	keySize := binary.BigEndian.Uint32(buf[crcSizeByteSize : crcSizeByteSize+keySizeByteSize])
	valueSize := binary.BigEndian.Uint64(buf[crcSizeByteSize+keySizeByteSize : crcSizeByteSize+keySizeByteSize+valueSizeByteSize])

	keyValueBuf := make([]byte, uint64(keySize)+valueSize)
	if _, err := io.ReadFull(r, keyValueBuf); err != nil {
		return nil, 0, err
	}

	var e *data.Entry = &data.Entry{}
	e.Crc = crc
	e.Key = types.Key(keyValueBuf[:keySize])
	e.Value = keyValueBuf[keySize : valueSize+uint64(keySize)]

	return e, uint64(crcSizeByteSize+keySizeByteSize+valueSizeByteSize) + uint64(keySize) + valueSize, nil
}

func NewBinaryCodec() Codec {
	return &binaryCodec{}
}
