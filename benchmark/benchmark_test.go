package benchmark

import (
	"fmt"
	"testing"

	"github.com/paolobroglio/kvstore/internal/storage"
)

func BenchmarkHashIndexWrite(b *testing.B) {
    idx := storage.NewHashIndex()
    store, err := storage.NewLogFile("./testdata", "bench.db", idx)
    if err != nil {
        b.Fatal(err)
    }
    defer store.Close()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        key := fmt.Sprintf("key_%d", i)
        value := fmt.Sprintf("value_%d", i)
        entry := &storage.Entry{
            Key:   []byte(key),
            Value: []byte(value),
        }
        store.Put(entry)
    }
}

func BenchmarkHashIndexRead(b *testing.B) {
}