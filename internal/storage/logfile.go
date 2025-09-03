package storage

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

type LogFile struct {
	writeHandle *os.File
	readHandle  *os.File
	mu          sync.Mutex
	offset      int64

	index Index
}

func NewLogFile(dbDir, dbFile string, index Index) (*LogFile, error) {
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, err
	}

	fullPath := filepath.Join(dbDir, dbFile)

	writeHandle, err := os.OpenFile(fullPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	readHandle, err := os.OpenFile(fullPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		writeHandle.Close()
		return nil, err
	}

	writeStat, err := writeHandle.Stat()
	if err != nil {
		writeHandle.Close()
		readHandle.Close()
		return nil, err
	}

	newLogFile := &LogFile{
		writeHandle: writeHandle,
		readHandle:  readHandle,
		offset:      writeStat.Size(),
		index:       index,
	}

	if writeStat.Size() > 0 {
		if err := newLogFile.rebuildIndex(); err != nil {
			newLogFile.Close()
			return nil, fmt.Errorf("failed to rebuild index: %w", err)
		}
	}

	return newLogFile, nil
}

func (lf *LogFile) rebuildIndex() error {
	_, err := lf.readHandle.Seek(0, 0)
	if err != nil {
		return err
	}

	offset := int64(0)
	reader := bufio.NewReader(lf.readHandle)

	for {
		currentOffset := offset
		lengths := make([]byte, 8)
		_, err := io.ReadFull(reader, lengths)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			break
		}
		if err != nil {
			return err
		}

		keyLen := binary.LittleEndian.Uint32(lengths[0:4])
		valueLen := binary.LittleEndian.Uint32(lengths[4:8])

		key := make([]byte, keyLen)
		keyBytesRead, err := io.ReadFull(reader, key)
		if err != nil {
			return err
		}

		value := make([]byte, valueLen)
		valueBytesRead, err := io.ReadFull(reader, value)
		if err != nil {
			return err
		}

		bytesRead := 8 + keyBytesRead + valueBytesRead

		location := Location{
			FileID: 0,
			Offset: currentOffset,
			Size:   int32(bytesRead),
		}

		if err := lf.index.Put(key, location); err != nil {
			return err
		}

		offset += int64(bytesRead)
	}

	lf.offset = offset
	
	return nil
}

func (lf *LogFile) Put(entry *Entry) error {
	lf.mu.Lock()
	defer lf.mu.Unlock()

	data, err := entry.Serialize()
	if err != nil {
		return err
	}

	currentOffset := lf.offset

	n, err := lf.writeHandle.Write(data)
	if err != nil {
		return err
	}

	lf.offset += int64(n)

	location := Location{
		FileID: 0,
		Offset: currentOffset,
		Size:   int32(n),
	}

	//log.Printf("Created new entry at location %+v\n", location)

	return lf.index.Put(entry.Key, location)
}

func (lf *LogFile) Get(key []byte) (*Entry, error) {

	location := lf.index.Get(key)
	if location == nil {
		return nil, nil
	}

	//log.Printf("Getting entry at location %+v\n", location)

	data := make([]byte, location.Size)
	_, err := lf.readHandle.ReadAt(data, location.Offset)
	if err != nil {
		return nil, err
	}

	return Deserialize(data)
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
