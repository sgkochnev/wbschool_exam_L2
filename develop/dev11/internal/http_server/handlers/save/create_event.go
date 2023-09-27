package save

import (
	"calendar/internal/http_server/handlers/utils"
	"calendar/internal/model"
	"log/slog"
	"net/http"
)

// Saver interface describes methods for saving events
type Saver interface {
	Save(event *model.Event) error
}

// Execute handler for saving events
func Execute(log *slog.Logger, store Saver) http.Handler {
	createEvent := func(w http.ResponseWriter, r *http.Request) {
		op := "save.CreateEvent"

		log = log.With(
			slog.String("op", op),
		)

		if r.Method != http.MethodPost {
			log.Error("method not allowed", "method", r.Method)
			utils.RenderJSON(w, http.StatusMethodNotAllowed, utils.Response{Message: "method not allowed"})
			return
		}

		//валидация
		event, err := utils.GetEvent(r)
		if err != nil {
			log.Error("parse: invalid request", "err", err.Error())
			utils.RenderJSON(w, http.StatusBadRequest, utils.Response{Message: "invalid request"})
			return
		}

		err = utils.ValidateEvent(event)
		if err != nil {
			log.Error("validate: invalid request", "err", err.Error())
			utils.RenderJSON(w, http.StatusBadRequest, utils.Response{Message: "invalid request"})
			return
		}

		//сохранение
		err = store.Save(event)
		if err != nil {
			log.Error("internal error", "err", err.Error())
			utils.RenderJSON(w, http.StatusInternalServerError, utils.Response{Message: "internal error"})
			return
		}

		//отправка ответа
		//w.WriteHeader(http.StatusCreated)
		utils.RenderJSON(w, http.StatusCreated, utils.Response{Message: "ok"})
		return
	}

	return http.HandlerFunc(createEvent)
}
