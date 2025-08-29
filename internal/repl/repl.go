package repl

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/paolobroglio/kvstore/internal/storage"
)

type REPL struct {
	store  storage.Storage
	scanner *bufio.Scanner
}

func New(store storage.Storage) *REPL {
	return &REPL{
		store:   store,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (r *REPL) Start() error {
	fmt.Println("Key-Value Database REPL")
	fmt.Println("Commands: put <key> <value>, get <key>, quit")
	
	for {
		fmt.Print("> ")
		
		if !r.scanner.Scan() {
			break
		}
		
		line := r.scanner.Text()
		if err := r.handleCommand(line); err != nil {
			if err == ErrQuit {
				break
			}
			fmt.Printf("Error: %v\n", err)
		}
	}
	
	return r.scanner.Err()
}

var ErrQuit = fmt.Errorf("quit command")

func (r *REPL) handleCommand(line string) error {
	tokens := strings.Fields(line)
	if len(tokens) == 0 {
		return nil
	}
	
	switch strings.ToLower(tokens[0]) {
	case "quit", "q":
		return ErrQuit
		
	case "put":
		if len(tokens) < 3 {
			fmt.Println("Usage: put <key> <value>")
			return nil
		}
		entry := &storage.Entry{
			Key:   []byte(tokens[1]),
			Value: []byte(tokens[2]),
		}
		return r.store.Put(entry)
		
	case "get":
		if len(tokens) < 2 {
			fmt.Println("Usage: get <key>")
			return nil
		}
		
		entry, err := r.store.Get([]byte(tokens[1]))
		if err != nil {
			return err
		}
		
		if entry == nil {
			fmt.Println("entry not found")
		} else {
			fmt.Printf("%q\n", string(entry.Value))
		}
		
	default:
		fmt.Println("Invalid command")
	}
	
	return nil
}