package storage

import (
	"encoding/binary"
	"io"
)

type Entry struct {
	Key   []byte
	Value []byte
}

func (e *Entry) Serialize() []byte {
	data := make([]byte, 8+len(e.Key)+len(e.Value))
	
	binary.LittleEndian.PutUint32(data[0:4], uint32(len(e.Key)))
	binary.LittleEndian.PutUint32(data[4:8], uint32(len(e.Value)))
	
	copy(data[8:8+len(e.Key)], e.Key)
	copy(data[8+len(e.Key):], e.Value)
	
	return data
}

func DeserializeEntry(r io.Reader) (*Entry, error) {
	lengths := make([]byte, 8)
	_, err := io.ReadFull(r, lengths)
	if err != nil {
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return nil, nil // End of file
		}
		return nil, err
	}

	keyLen := binary.LittleEndian.Uint32(lengths[0:4])
	valueLen := binary.LittleEndian.Uint32(lengths[4:8])

	key := make([]byte, keyLen)
	if _, err := io.ReadFull(r, key); err != nil {
		return nil, err
	}

	value := make([]byte, valueLen)
	if _, err := io.ReadFull(r, value); err != nil {
		return nil, err
	}

	return &Entry{Key: key, Value: value}, nil
}