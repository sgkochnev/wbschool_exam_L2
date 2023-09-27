package storage

import "calendar/internal/model"

// Update - update event in storage
func (s *Storage) Update(event *model.Event) error {
	path := s.getPath(event.UserId, event.Date)

	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.getStore(path)
	if err != nil {
		return err
	}
	if _, ok := store[event.Date]; !ok {
		return ErrNotFound
	}

	store[event.Date] = *event

	return s.commit(path, store)
}
