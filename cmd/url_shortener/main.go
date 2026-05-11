package main

import (
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/setup"
	"url-shortener/internal/lib/logger/sl"
	redirect "url-shortener/internal/server/handlers/redirect"
	"url-shortener/internal/server/handlers/url/delete"
	"url-shortener/internal/server/handlers/url/save"
	"url-shortener/internal/server/middleware/logger"
	"url-shortener/internal/storage/postgres"
	"url-shortener/internal/storage/redis"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.Mustload()

	log := setup.SetupLogger(cfg.Env)

	log.Info("starting url-shotener service", slog.String("env", cfg.Env))
	log.Debug("debug message are enabled")

	storage, err := postgres.New(&cfg.Storage)
	if err != nil {
		log.Error("Cannot open storage", sl.Err(err))
		os.Exit(1)
	}

	cache, err := redis.New(&cfg.Cache)
	if err != nil {
		log.Error("Cannot open cache", sl.Err(err))
		os.Exit(1)
	}

	// Can do with base router net/http
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.Server.User: cfg.Server.Password,
		}))

		r.Post("/", save.New(log, storage))
		r.Delete("/{alias}", delete.New(log, storage))
	})

	router.Get("/{alias}", redirect.New(log, storage, cache))

	log.Info("starting url-shotener server", slog.String("address", cfg.Server.Address))

	srv := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed with start server")
	}

	log.Error("server stopped")
}
