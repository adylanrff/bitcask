package codec

import (
	"io"

	"github.com/adylanrff/bitcask/internal/data"
)

type Codec interface {
	Encode(w io.Writer, entry *data.Entry) (uint64, error)
	Decode(r io.Reader) (*data.Entry, uint64, error)
}
