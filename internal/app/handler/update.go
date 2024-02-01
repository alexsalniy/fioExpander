package handler

import (
	"fio-expander/internal/app/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func Update(log *slog.Logger, fioexp FIOExp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.ExtendedFIO

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)
			render.JSON(w, r, Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("requset", req))

		err = fioexp.Update(&req)
		if err != nil {
			log.Error("failed to update fio in database", err)

			render.JSON(w, r, Error("failed to update fio"))

			return
		}

		log.Info("fio added", slog.Any("id", req.ID))

		render.JSON(w, r, IDResponse{
			Response: OK(),
			ID:       req.ID.String(),
		})
	}
}
