# kvstore

A log-structured key-value store written in Go with pluggable indexing strategies.

## Usage

### Prerequisites
  
Go 1.19 or later

### Building

```bash
make build
```

### Running

```bash
# Start the KV store with hash index
make run

# Or specify index type
./bin/kvstore -index=hash
```

## Project status

**Under Development** - Currently implementing different indexing strategies

### Completed features

- **Append-only log**
  - Append-only log-structured storage
  - Binary serialization: key_length + value_length + key + value
- **Indexing**
  - Pluggable indexing strategy
    - HashIndex
  - Index rebuilding starting from the log file
- **Core Storage Functionalities**
  - Put
  - Get
  - Close
  - Key overwriting, latest value wins
- **REPL**
  - Get/Put commands

### In Progress

- **Benchmarking framework**
  - custom benchmarking for performance comparison
  - CPU and Memory profiling

### Planned Features

- **Additional Indexing Strategies**
  - B-Tree
  - LSM-Tree
- **Advanced storage functionalities**
  - Log file rotation and compaction
  - Parallelization
- **REPL improvements**
  - more robust parser
  - DELETE command
- **Config**
  - more configuration options
- **TBD: HTTP Client/Server**


