package update

import (
	"calendar/internal/http_server/handlers/utils"
	"calendar/internal/model"
	"log/slog"
	"net/http"
)

// Updater interface describes methods for updating events
type Updater interface {
	Update(event *model.Event) error
}

// Execute handler for updating events
func Execute(log *slog.Logger, store Updater) http.Handler {
	updateEvent := func(w http.ResponseWriter, r *http.Request) {
		op := "update.UpdateEvent"

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
		err = store.Update(event)
		if err != nil {
			log.Error("internal error", "err", err.Error())
			utils.RenderJSON(w, http.StatusServiceUnavailable, utils.Response{Message: "internal error"})
			return
		}

		//отправка ответа
		utils.RenderJSON(w, http.StatusOK, utils.Response{Message: "ok"})
	}

	return http.HandlerFunc(updateEvent)
}
