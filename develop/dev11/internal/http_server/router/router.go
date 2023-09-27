package router

import (
	"calendar/internal/http_server/handlers/delete"
	"calendar/internal/http_server/handlers/get"
	"calendar/internal/http_server/handlers/save"
	"calendar/internal/http_server/handlers/update"
	l "calendar/pkg/logger"

	"calendar/internal/http_server/middleware/logger"
	"calendar/internal/model"
	"log/slog"
	"net/http"
)

// Storage interface describes methods for interacting with storage
type Storage interface {
	Save(event *model.Event) error
	Delete(userId int, date string) error
	Update(event *model.Event) error
	EventsForDay(userId int, date string) (*model.Event, error)
	EventsForWeek(userId int, date string) ([]model.Event, error)
	EventsForMonth(userId int, date string) ([]model.Event, error)
}

// New create new router
func New(store Storage) *http.ServeMux {
	router := http.NewServeMux()

	log := l.New(slog.LevelDebug)
	loggingMiddleware := logger.New(log)

	router.Handle("/create_event", loggingMiddleware(save.Execute(log, store)))
	router.Handle("/update_event", loggingMiddleware(update.Execute(log, store)))
	router.Handle("/delete_event", loggingMiddleware(delete.Execute(log, store)))
	router.Handle("/events_for_day", loggingMiddleware(get.EventsForDay(log, store)))
	router.Handle("/events_for_week", loggingMiddleware(get.EventsForWeek(log, store)))
	router.Handle("/events_for_month", loggingMiddleware(get.EventsForMonth(log, store)))

	return router
}
