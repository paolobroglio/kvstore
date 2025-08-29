package storage

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
)

type LogFile struct {
	writeHandle *os.File
	readHandle  *os.File
}

func NewLogFile(dbDir, dbFile string) (*LogFile, error) {
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(dbDir, "db.txt")

	writeHandle, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	readHandle, err := os.OpenFile(fullPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		writeHandle.Close()
		return nil, err
	}

	return &LogFile{
		writeHandle: writeHandle,
		readHandle:  readHandle,
	}, nil
}

func (lf *LogFile) Put(entry *Entry) error {
	data := entry.Serialize()
	_, err := lf.writeHandle.Write(data)
	return err
}

func (lf *LogFile) Get(key []byte) (*Entry, error) {
	if _, err := lf.readHandle.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	for {
		entry, err := DeserializeEntry(lf.readHandle)
		if err != nil {
			return nil, err
		}
		if entry == nil {
			break
		}

		if bytes.Equal(entry.Key, key) {
			return entry, nil
		}
	}

	return nil, nil
}

func (lf *LogFile) Close() error {
	var err1, err2 error
	if lf.writeHandle != nil {
		err1 = lf.writeHandle.Close()
	}
	if lf.readHandle != nil {
		err2 = lf.readHandle.Close()
	}

	if err1 != nil {
		return err1
	}
	return err2
}
