package config

import "os"

type Config struct {
	DBDir  string
	DBFile string
}

func New() *Config {
	dbDir := os.Getenv("KVSTORE_DIR")
	if dbDir == "" {
		dbDir = "db"
	}
	
	return &Config{
		DBDir:  dbDir,
		DBFile: "db.txt",
	}
}