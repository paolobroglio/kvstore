package benchmark

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sort"
	"time"

	"github.com/paolobroglio/kvstore/internal/storage"
)

type Config struct {
	IndexType  string
	NumKeys    int
	ValueSize  int
	ReadRatio  float64
	Duration   time.Duration
	KeyPattern string
}

type Result struct {
	IndexType       string
	ThroughputOps   float64
	AvgReadLatency  time.Duration
	AvgWriteLatency time.Duration
	P95ReadLatency  time.Duration
	P99ReadLatency  time.Duration
	MemoryUsageMB   float64
	StartupTimeMs   int64
}

func Run(config Config) (*Result, error) {
	var idx storage.Index
	switch config.IndexType {
	case "hash":
		idx = storage.NewHashIndex()
	default:
		return nil, fmt.Errorf("unknown index type: %s", config.IndexType)
	}
	return runBenchmark(idx, config)
}

func runBenchmark(idx storage.Index, config Config) (*Result, error) {
	benchDir := fmt.Sprintf("./benchmark/testdata_%d", time.Now().UnixNano())
	if err := os.MkdirAll(benchDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create benchmark directory: %w", err)
	}
	defer os.RemoveAll(benchDir)

	var memBefore runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memBefore)

	startupStart := time.Now()
	store, err := storage.NewLogFile(benchDir, "bench.db", idx)
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}
	defer store.Close()
	startupTime := time.Since(startupStart)

	keys := generateKeys(config.NumKeys, config.KeyPattern)
	values := generateValues(config.NumKeys, config.ValueSize)

	populateCount := int(float64(config.NumKeys) * 0.8)
	fmt.Printf("Pre-populating %d entries...\n", populateCount)

	for i := 0; i < populateCount; i++ {
		entry := &storage.Entry{
			Key:   []byte(keys[i]),
			Value: []byte(values[i]),
		}
		if err := store.Put(entry); err != nil {
			return nil, fmt.Errorf("failed to pre-populate data: %w", err)
		}
	}

	var memAfterPopulate runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memAfterPopulate)

	var readLatencies []time.Duration
	var writeLatencies []time.Duration

	rand.Seed(time.Now().UnixNano())

	fmt.Printf("Running workload for %v...\n", config.Duration)
	start := time.Now()
	operations := 0
	readOps := 0
	writeOps := 0

	for time.Since(start) < config.Duration {
		if rand.Float64() < config.ReadRatio {
			keyIdx := rand.Intn(populateCount)
			key := []byte(keys[keyIdx])

			opStart := time.Now()
			_, err := store.Get(key)
			latency := time.Since(opStart)

			readLatencies = append(readLatencies, latency)
			readOps++

			if err != nil && err.Error() != "key not found" {
				return nil, fmt.Errorf("read operation failed: %w", err)
			}
		} else {
			keyIdx := rand.Intn(config.NumKeys)
			entry := &storage.Entry{
				Key:   []byte(keys[keyIdx]),
				Value: []byte(generateRandomValue(config.ValueSize)),
			}

			opStart := time.Now()
			err := store.Put(entry)
			latency := time.Since(opStart)

			writeLatencies = append(writeLatencies, latency)
			writeOps++

			if err != nil {
				return nil, fmt.Errorf("write operation failed: %w", err)
			}
		}
		operations++

		if operations%1000 == 0 {
			elapsed := time.Since(start)
			fmt.Printf("  %d ops in %v (%.1f ops/sec)\n",
				operations, elapsed, float64(operations)/elapsed.Seconds())
		}
	}

	totalDuration := time.Since(start)

	var memAfter runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memAfter)

	result := &Result{
		IndexType:     config.IndexType,
		ThroughputOps: float64(operations) / totalDuration.Seconds(),
		StartupTimeMs: startupTime.Milliseconds(),
		MemoryUsageMB: float64(memAfter.Alloc-memBefore.Alloc) / 1024 / 1024,
	}

	if len(readLatencies) > 0 {
		result.AvgReadLatency = average(readLatencies)
		result.P95ReadLatency = percentile(readLatencies, 0.95)
		result.P99ReadLatency = percentile(readLatencies, 0.99)
	}

	if len(writeLatencies) > 0 {
		result.AvgWriteLatency = average(writeLatencies)
	}

	fmt.Printf("\nBenchmark completed!\n")
	fmt.Printf("Total operations: %d (%d reads, %d writes)\n", operations, readOps, writeOps)
	fmt.Printf("Duration: %v\n", totalDuration)
	fmt.Printf("Throughput: %.1f ops/sec\n", result.ThroughputOps)

	return result, nil
}

func generateKeys(count int, pattern string) []string {
	keys := make([]string, count)

	switch pattern {
	case "sequential":
		for i := 0; i < count; i++ {
			keys[i] = fmt.Sprintf("key_%06d", i)
		}
	case "random":
		for i := 0; i < count; i++ {
			keys[i] = fmt.Sprintf("key_%06d", rand.Intn(count*10))
		}
	default:
		for i := 0; i < count; i++ {
			keys[i] = fmt.Sprintf("key_%06d", i)
		}
	}

	return keys
}

func generateValues(count, size int) []string {
	values := make([]string, count)
	for i := 0; i < count; i++ {
		values[i] = generateRandomValue(size)
	}
	return values
}

func generateRandomValue(size int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, size)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func average(durations []time.Duration) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	var total time.Duration
	for _, d := range durations {
		total += d
	}
	return total / time.Duration(len(durations))
}

func percentile(durations []time.Duration, p float64) time.Duration {
	if len(durations) == 0 {
		return 0
	}

	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	index := int(math.Ceil(p*float64(len(sorted)))) - 1
	if index < 0 {
		index = 0
	}
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}
