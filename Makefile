.PHONY: build test clean run

build:
	go build -o bin/kvstore cmd/kvstore/main.go

test:
	go test -v ./...

clean:
	rm -rf bin/ db/

run: build
	./bin/kvstore