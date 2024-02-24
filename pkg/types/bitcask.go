package types

type Key string
type Value []byte

type FoldFunc[A any] func(Key, Value, A) A

type BitcaskHandler interface {
	Get(key Key) (Value, error)
	Put(key Key, value Value) error
	Delete(key Key) error
	ListKeys() ([]Key, error)
	Fold(f FoldFunc[any], acc0 any) (any, error)
	Sync() error
	Close() error
}
