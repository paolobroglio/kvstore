package storage

type Storage interface {
	Put(entry *Entry) error
	Get(key []byte) (*Entry, error)
	Close() error
}