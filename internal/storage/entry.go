package storage

import (
	"encoding/binary"
	"fmt"
)

type Entry struct {
	Key   []byte
	Value []byte
}

func (e *Entry) Serialize() ([]byte, error) {
	keyLen := len(e.Key)
	valueLen := len(e.Value)

	totalSize := 8 + keyLen + valueLen
	buf := make([]byte, totalSize)

	binary.LittleEndian.PutUint32(buf[0:4], uint32(keyLen))
	binary.LittleEndian.PutUint32(buf[4:8], uint32(valueLen))
	
	copy(buf[8:8+keyLen], e.Key)
	copy(buf[8+keyLen:], e.Value)

	return buf, nil
}

func Deserialize(data []byte) (*Entry, error) {
	keyLen := binary.LittleEndian.Uint32(data[0:4])
	valueLen := binary.LittleEndian.Uint32(data[4:8])

	expectedSize := 8 + int(keyLen) + int(valueLen)
	if len(data) < expectedSize {
		return nil, fmt.Errorf("data too short for entry content")
	}

	key := make([]byte, keyLen)
	value := make([]byte, valueLen)

	copy(key, data[8:8+keyLen])
	copy(value, data[8+keyLen:8+keyLen+valueLen])

	return &Entry{Key: key, Value: value}, nil
}
