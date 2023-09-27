package storage

import (
	"fmt"
	"sync"
)

var ErrNotFound = fmt.Errorf("not found")
var ErrAlreadyExists = fmt.Errorf("already exists")

// Storage - storage
type Storage struct {
	mu   sync.Mutex
	path string
}

// New - create storage
func New(path string) *Storage {
	return &Storage{
		path: path,
	}
}
