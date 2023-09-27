package storage

import (
	"calendar/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func (s *Storage) getStore(f string) (map[string]model.Event, error) {
	data, err := os.ReadFile(f)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	var events map[string]model.Event
	if len(data) == 0 {
		events = make(map[string]model.Event)
		return events, nil
	}

	err = json.Unmarshal(data, &events)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (s *Storage) commit(f string, events map[string]model.Event) error {
	data, err := json.Marshal(events)
	if err != nil {
		return err
	}
	return os.WriteFile(f, data, os.ModePerm)
}

func (s *Storage) getPath(userId int, date string) string {
	d, _ := time.Parse("2006-01-02", date)
	return filepath.Join(s.path,
		fmt.Sprintf("%d", userId),
		fmt.Sprintf("%d", d.Year()),
		fmt.Sprintf("%d.json", d.Month()),
	)
}
