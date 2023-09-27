package storage

import (
	"calendar/internal/model"
	"fmt"
	"slices"
	"time"
)

// EventsForDay - get events for day
func (s *Storage) EventsForDay(userId int, date string) (*model.Event, error) {
	path := s.getPath(userId, date)

	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.getStore(path)
	if err != nil {
		return nil, err
	}

	event, ok := store[date]
	if !ok {
		return nil, nil
	}

	return &event, nil
}

// EventsForWeek - get events for week
func (s *Storage) EventsForWeek(userId int, date string) ([]model.Event, error) {
	path := s.getPath(userId, date)

	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.getStore(path)
	if err != nil {
		return nil, err
	}

	week := [7]string{0: date}
	d1, _ := time.Parse("2006-01-02", date)
	d2 := d1
	for i := 1; i < 7; i++ {
		d2 = d2.AddDate(0, 0, 1)
		week[i] = d2.Format("2006-01-02")
	}

	var events []model.Event

	if d1.Month() != d2.Month() {
		path = s.getPath(userId, week[6])
		store, err := s.getStore(path)
		if err != nil {
			return nil, err
		}
		for _, day := range week {
			if _, ok := store[day]; ok {
				events = append(events, store[day])
			}
		}
	}

	for _, day := range week {
		if _, ok := store[day]; ok {
			events = append(events, store[day])
		}
	}

	slices.SortFunc(events, func(a, b model.Event) int {
		switch {
		case a.Date < b.Date:
			return -1
		case a.Date > b.Date:
			return 1
		default:
			return 0
		}
	})

	return events, nil
}

// EventsForMonth - get events for month
func (s *Storage) EventsForMonth(userId int, date string) ([]model.Event, error) {
	path := s.getPath(userId, date)
	fmt.Println(path)
	s.mu.Lock()
	defer s.mu.Unlock()

	store, err := s.getStore(path)
	if err != nil {
		return nil, err
	}

	var events []model.Event
	for _, event := range store {
		events = append(events, event)
	}

	slices.SortFunc(events, func(a, b model.Event) int {
		switch {
		case a.Date < b.Date:
			return -1
		case a.Date > b.Date:
			return 1
		default:
			return 0
		}
	})
	fmt.Println(events)
	return events, nil
}
