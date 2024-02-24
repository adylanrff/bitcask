package minimal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/adylanrff/bitcask/internal/codec"
	"github.com/adylanrff/bitcask/internal/data"
	"github.com/adylanrff/bitcask/pkg/option"
	"github.com/adylanrff/bitcask/pkg/types"
)

// Minimal bitcask implements the bitcask handler with bare minimum capabilities.
type Handler struct {
	sync.RWMutex
	dirname string
	opts    *option.Options

	inmem      data.Keydir
	activeFile *os.File
	offset     uint64

	codec codec.Codec
}

func NewHandler(dirname string, options *option.Options) (*Handler, error) {

	h := &Handler{
		dirname: dirname,
		opts:    options,
		codec:   codec.NewBinaryCodec(),
		inmem:   make(data.Keydir),
	}

	return h, nil
}

func (h *Handler) Init() error {
	if _, err := os.Stat(h.dirname); os.IsNotExist(err) {
		// Directory does not exist, create it
		err := os.Mkdir(h.dirname, 0755) // 0755 is the permission mode
		if err != nil {
			return err
		}
	}

	if err := h.createActiveFileIfNotExist(); err != nil {
		return err
	}

	if err := h.buildInmemKeydir(); err != nil {
		return err
	}

	return nil
}

// Get implements bitcask.Handler.
func (h *Handler) Get(key types.Key) (types.Value, error) {
	keydir, ok := h.inmem.Get(key)
	if !ok {
		return nil, nil
	}

	value, err := h.readValue(keydir)
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Put implements bitcask.Handler.
func (h *Handler) Put(key types.Key, value types.Value) error {
	entry := data.NewEntry(
		uint64(time.Now().UnixMilli()),
		key,
		value,
	)

	return h.appendEntry(entry)
}

// Delete implements bitcask.Handler.
func (*Handler) Delete(key types.Key) error {
	panic("unimplemented")
}

// Fold implements bitcask.Handler.
func (*Handler) Fold(f types.FoldFunc[any], acc0 any) (any, error) {
	panic("unimplemented")
}

// ListKeys implements bitcask.Handler.
func (*Handler) ListKeys() ([]types.Key, error) {
	panic("unimplemented")
}

// Sync implements bitcask.Handler.
func (*Handler) Sync() error {
	panic("unimplemented")
}

// Close implements bitcask.Handler.
func (h *Handler) Close() error {
	return h.Sync()
}

func (h *Handler) buildInmemKeydir() error {
	dirEntries, err := os.ReadDir(h.dirname)
	if err != nil {
		return err
	}

	dirFiles := make([]string, 0, len(dirEntries))
	for _, dirEntry := range dirEntries {
		dirFiles = append(dirFiles, dirEntry.Name())
	}
	sort.Strings(dirFiles)

	dataOffset := uint64(0)

	for _, filename := range dirFiles {
		f, err := os.Open(filepath.Join(h.dirname, filename))
		if err != nil {
			return err
		}

		for true {
			decodedData, n, err := h.codec.Decode(f)
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}

			dataOffset += n

			h.inmem.Put(decodedData.Key, data.KeydirEntry{
				FileID:      f.Name(),
				ValueSize:   uint64(len(decodedData.Value)),
				ValueOffset: uint64(dataOffset + uint64(16) + uint64(len(decodedData.Key))),
				Timestamp:   decodedData.Timestamp,
			})
		}
	}

	return nil
}

func (h *Handler) appendEntry(entry *data.Entry) error {

	h.Lock()
	defer h.Unlock()
	// START critical section
	n, err := h.codec.Encode(h.activeFile, entry)
	if err != nil {
		return err
	}

	h.inmem.Put(entry.Key, data.KeydirEntry{
		FileID:      h.activeFile.Name(),
		ValueSize:   uint64(len(entry.Value)),
		ValueOffset: h.offset + uint64(16) + uint64(len(entry.Key)),
		Timestamp:   entry.Timestamp,
	})

	h.offset += n
	// END Critical Section

	return nil
}

func (h *Handler) readEntry(keydirEntry data.KeydirEntry) (*data.Entry, error) {
	filename := keydirEntry.FileID
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	if _, err := f.Seek(int64(keydirEntry.ValueOffset), 0); err != nil {
		return nil, err
	}

	decodedData, _, err := h.codec.Decode(f)
	if err != nil {
		return nil, err
	}

	return decodedData, nil
}

func (h *Handler) readValue(keydirEntry data.KeydirEntry) (types.Value, error) {
	filename := keydirEntry.FileID
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	buf := make([]byte, keydirEntry.ValueSize)

	if _, err := f.Seek(int64(keydirEntry.ValueOffset), 0); err != nil {
		return nil, err
	}

	if _, err := io.ReadFull(f, buf); err != nil {
		return nil, err
	}

	return types.Value(buf), nil
}

func (h *Handler) createActiveFileIfNotExist() error {
	if h.activeFile != nil {
		return nil
	}

	filename := fmt.Sprintf("%s/%s-%d.data", h.dirname, h.dirname, time.Now().UnixMilli())
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	h.activeFile = file
	return nil
}
