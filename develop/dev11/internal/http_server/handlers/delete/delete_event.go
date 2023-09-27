package delete

import (
	"calendar/internal/http_server/handlers/utils"
	"fmt"
	"log/slog"
	"net/http"
)

// Deleter interface describes methods for deleting events
type Deleter interface {
	Delete(userId int, date string) error
}

// Execute handler for deleting events
func Execute(log *slog.Logger, store Deleter) http.Handler {
	deleteEvent := func(w http.ResponseWriter, r *http.Request) {
		op := "delete.DeleteEvent"
		log = log.With(
			slog.String("op", op),
		)

		if r.Method != http.MethodPost {
			log.Error("method not allowed", "method", r.Method)
			utils.RenderJSON(w, http.StatusMethodNotAllowed, utils.Response{Message: "method not allowed"})
			return
		}

		//валидация
		dto, err := utils.GetDTO(r)
		if err != nil {
			log.Error("parse: invalid request", "err", err.Error())
			utils.RenderJSON(w, http.StatusBadRequest, utils.Response{Message: "invalid request"})
			return
		}
		err = utils.ValidateDTO(dto)
		if err != nil {
			log.Error("validate: invalid request", "err", err.Error())
			utils.RenderJSON(w, http.StatusBadRequest, utils.Response{Message: "invalid request"})
			return
		}
		fmt.Println(dto)

		//удаление
		err = store.Delete(dto.UserId, dto.Date)
		if err != nil {
			log.Error("internal error", "err", err.Error())
			utils.RenderJSON(w, http.StatusInternalServerError, utils.Response{Message: "internal error"})
			return
		}

		//отправка ответа
		utils.RenderJSON(w, http.StatusOK, utils.Response{Message: "ok"})
	}

	return http.HandlerFunc(deleteEvent)
}
