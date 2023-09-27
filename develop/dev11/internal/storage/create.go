package storage

import (
	"calendar/internal/model"
	"os"
	"path/filepath"
)

// Save - save event to storage
func (s *Storage) Save(event *model.Event) error {

	path := s.getPath(event.UserId, event.Date)

	s.mu.Lock()
	defer s.mu.Unlock()

	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return err
	}

	store, err := s.getStore(path)
	if err != nil {
		return err
	}
	if _, ok := store[event.Date]; ok {
		return ErrAlreadyExists
	}
	store[event.Date] = *event

	return s.commit(path, store)
}
