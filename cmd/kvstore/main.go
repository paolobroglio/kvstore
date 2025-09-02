package main

import (
	"log"

	"github.com/paolobroglio/kvstore/internal/config"
	"github.com/paolobroglio/kvstore/internal/repl"
	"github.com/paolobroglio/kvstore/internal/storage"
)

func main() {
	cfg := config.New()
	
	store, err := storage.NewLogFile(cfg.DBDir, cfg.DBFile, storage.NewHashIndex())
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	defer store.Close()

	r := repl.New(store)
	if err := r.Start(); err != nil {
		log.Fatalf("REPL error: %v", err)
	}
}