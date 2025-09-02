package main

import (
	"fmt"
	"log"
	"time"

	"github.com/paolobroglio/kvstore/benchmark"
)

func main() {
	config := benchmark.Config{
		IndexType:  "hash",
		NumKeys:    10000,
		ValueSize:  100,
		ReadRatio:  0.8,
		Duration:   30 * time.Second,
		KeyPattern: "random",
	}

	fmt.Printf("Running benchmark with config: %+v\n", config)

	result, err := benchmark.Run(config)
	if err != nil {
		log.Fatal(err)
	}

	printResults(result)
}

func printResults(r *benchmark.Result) {
	fmt.Printf("\n=== Benchmark Results ===\n")
	fmt.Printf("Index Type:      %s\n", r.IndexType)
	fmt.Printf("Throughput:      %.1f ops/sec\n", r.ThroughputOps)
	fmt.Printf("Avg Read:        %v\n", r.AvgReadLatency)
	fmt.Printf("Avg Write:       %v\n", r.AvgWriteLatency)
	fmt.Printf("P95 Read:        %v\n", r.P95ReadLatency)
	fmt.Printf("P99 Read:        %v\n", r.P99ReadLatency)
	fmt.Printf("Memory Usage:    %.1f MB\n", r.MemoryUsageMB)
	fmt.Printf("Startup Time:    %d ms\n", r.StartupTimeMs)
}
