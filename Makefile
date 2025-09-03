.PHONY: build test clean run bench bench-run bench-profile bench-compare bench-clean

build:
	go build -o bin/kvstore cmd/kvstore/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/ db/

run: build
	./bin/kvstore

bench:
	go test -bench=. -benchmem ./benchmark/

bench-run:
	go run ./benchmark/cmd/ -index=hash -keys=10000 -duration=30s

bench-profile:
	go test -bench=. -cpuprofile=cpu.prof -memprofile=mem.prof ./benchmark/
	@echo "Analyze profiles with:"
	@echo "  go tool pprof cpu.prof"
	@echo "  go tool pprof mem.prof"

# bench-compare:
# 	@echo "=== Hash Index Benchmark ==="
# 	go run ./benchmark/cmd/ -index=hash -keys=10000 -duration=10s
# 	@echo ""
# 	@echo "=== B-Tree Index Benchmark ==="
# 	go run ./benchmark/cmd/ -index=btree -keys=10000 -duration=10s

bench-clean:
	rm -rf ./benchmark/testdata/
	rm -f *.prof
	rm -f ./benchmark/*.prof

bench-main:
	go run ./cmd/kvstore/ -benchmark -index=hash

bench-suite:
	@echo "Running comprehensive benchmark suite..."
	@echo "=== Small dataset, read-heavy ==="
	go run ./benchmark/cmd/ -index=hash -keys=1000 -read-ratio=0.9 -duration=10s
	@echo ""
	@echo "=== Medium dataset, balanced ==="
	go run ./benchmark/cmd/ -index=hash -keys=10000 -read-ratio=0.5 -duration=15s
	@echo ""
	@echo "=== Large dataset, write-heavy ==="
	go run ./benchmark/cmd/ -index=hash -keys=50000 -read-ratio=0.1 -duration=20s