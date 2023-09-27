package storage

// Delete - delete event from storage
func (s *Storage) Delete(userId int, date string) error {
	path := s.getPath(userId, date)

	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.getStore(path)
	if err != nil {
		return err
	}
	delete(store, date)

	return s.commit(path, store)
}
