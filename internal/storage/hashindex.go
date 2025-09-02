package storage

import "sync"

type HashIndex struct {
	mu sync.RWMutex
	index map[string]Location
}

func NewHashIndex() *HashIndex {
	return &HashIndex{
		index: make(map[string]Location),
	}
}

func (hi *HashIndex) Put(key []byte, location Location) error {
	hi.mu.Lock()
	defer hi.mu.Unlock()

	hi.index[(string(key))] = location
	return nil
}
func (hi *HashIndex) Get(key []byte) *Location {
	hi.mu.RLock()
	defer hi.mu.RUnlock()

	if loc, exists := hi.index[string(key)]; exists {
    return &loc
  }
	return nil
}
func (hi *HashIndex) Delete(key []byte) error {
	hi.mu.Lock()
  defer hi.mu.Unlock()
    
  delete(hi.index, string(key))
  return nil
}
func (hi *HashIndex) Close() error {
	return nil
}