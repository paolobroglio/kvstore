package config

import "os"

type Config struct {
	DBDir     string
	DBFile    string
	IndexType string
}

func New() *Config {
	dbDir := os.Getenv("KVSTORE_DIR")
	if dbDir == "" {
		dbDir = "db"
	}

	indexType := os.Getenv("INDEX_TYPE")
	if indexType == "" {
		indexType = "hash"
	}

	return &Config{
		DBDir:  dbDir,
		DBFile: "store.db",
		IndexType: indexType,
	}
}
