package main

import (
	"fio-expander/internal/app/handler"
	"fio-expander/internal/config"
	"fio-expander/internal/storage/psql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	cfg := config.MustLoad()

	fmt.Println(cfg)

	var log *slog.Logger
	log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	log.Info("start fio-expander")

	storage, err := psql.New(cfg.Dbsource)
	if err != nil {
		log.Error("failed to init storage", err)
		os.Exit(1)
	}

	log.Info("database connection is successful")
	_ = storage

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/fio", handler.New(log, storage))
	router.Get("/fio", handler.FindBy(log, storage))
	router.Delete("/fio", handler.Delete(log, storage))
	router.Put("/fio", handler.Update(log, storage))

	log.Info("strting server", slog.String("address", cfg.Server))

	srv := &http.Server{
		Addr:    cfg.Server,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	log.Error("server stoped")

}
