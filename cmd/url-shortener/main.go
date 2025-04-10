package main

import (
	"github.com/UncleWeb1992/Url-SHortner/internal/config"
	"github.com/UncleWeb1992/Url-SHortner/internal/http-server/handlers/redirect"
	"github.com/UncleWeb1992/Url-SHortner/internal/http-server/handlers/url/delete"
	"github.com/UncleWeb1992/Url-SHortner/internal/http-server/handlers/url/save"
	"github.com/UncleWeb1992/Url-SHortner/internal/http-server/middleware/logger"
	"github.com/UncleWeb1992/Url-SHortner/internal/lib/logger/sl"
	"github.com/UncleWeb1992/Url-SHortner/internal/storage/sqlite"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// TODO: init config: cleanenv
	cfg := config.MustLoad()

	// TODO: init logger: slog
	log := setupLogger(cfg.Env)

	log.Info("starting url-shortener", slog.String("env", cfg.Env))

	// TODO: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("cannot init storage", sl.Err(err))
		os.Exit(1)
	}

	// TODO: init router: chi, "chi render"
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/api/save", save.New(log, storage, cfg.AliasLength))
	router.Get("/api/redirect/{alias}", redirect.New(log, storage))
	router.Delete("/api/{alias}", delete.New(log, storage))

	log.Info("server started", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}

	// TODO: start server
	if err := srv.ListenAndServe(); err != nil {
		log.Error("cannot start server", sl.Err(err))
	}
	log.Error("server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
