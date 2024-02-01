package handler

import (
	"fio-expander/internal/app/model"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
)

func Delete(log *slog.Logger, fioexp FIOExp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.ExtendedFIO

		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode request body", err)
			render.JSON(w, r, Error("failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("requset", req))

		err = fioexp.Delete(&req)
		if err != nil {
			log.Error("failed to delete fio from database", err)

			render.JSON(w, r, Error("failed to delete fio"))

			return
		}

		log.Info("fio added", slog.Any("id", req.ID))

		render.JSON(w, r, OK())

	}
}
