package storage

type Index interface {
	Put(key []byte, location Location) error
	Get(key []byte) *Location
	Delete(key []byte) error
	Close() error
}

type Location struct {
	FileID int64
	Offset int64
	Size int32
}