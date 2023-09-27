package get

import (
	"calendar/internal/http_server/handlers/utils"
	"calendar/internal/model"
	"log/slog"
	"net/http"
)

// GetterForDay interface describes methods for getting events for month
type GetterForMonth interface {
	EventsForMonth(userId int, date string) ([]model.Event, error)
}

// EventsForMonth handler for getting events for month
func EventsForMonth(log *slog.Logger, store GetterForMonth) http.Handler {
	eventsForMonth := func(w http.ResponseWriter, r *http.Request) {
		op := "get.EventsForMonth"

		log = log.With(
			slog.String("op", op),
		)

		var dto *utils.DTO
		var err error

		switch r.Method {
		case http.MethodGet:
			dto, err = utils.GetDTOFormParams(r)
			if err != nil {
				log.Error("parse: invalid request", "err", err.Error())
				utils.RenderJSON(w, http.StatusBadRequest, utils.Response{Message: "invalid request"})
				return
			}
		case http.MethodPost:
			dto, err = utils.GetDTO(r)
			if err != nil {
				log.Error("parse: invalid request", "err", err.Error())
				utils.RenderJSON(w, http.StatusBadRequest, utils.Response{Message: "invalid request"})
				return
			}
		default:
			utils.RenderJSON(w, http.StatusMethodNotAllowed, utils.Response{Message: "method not allowed"})
			return
		}

		err = utils.ValidateDTO(dto)
		if err != nil {
			log.Error("validate: invalid request", "err", err.Error())
			utils.RenderJSON(w, http.StatusBadRequest, utils.Response{Message: "invalid request"})
			return
		}

		event, err := store.EventsForMonth(dto.UserId, dto.Date)
		if err != nil {
			log.Error("internal error", "err", err.Error())
			utils.RenderJSON(w, http.StatusInternalServerError, utils.Response{Message: "internal error"})
			return
		}
		if event == nil {
			log.Info("event not found")
			utils.RenderJSON(w, http.StatusNotFound, utils.Response{Message: "event not found"})
			return
		}

		utils.RenderJSON(w, http.StatusOK, event)
	}

	return http.HandlerFunc(eventsForMonth)
}
