package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/clients"
	ssogrpc "url-shortener/internal/clients/sso/grpc"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/setup"
	"url-shortener/internal/lib/logger/sl"
	redirect "url-shortener/internal/server/handlers/redirect"
	"url-shortener/internal/server/handlers/url/delete"
	"url-shortener/internal/server/handlers/url/save"
	"url-shortener/internal/server/middleware/logger"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	log := setup.SetupLogger(cfg.Env)

	log.Info("starting url-shotener service", slog.String("env", cfg.Env))
	log.Debug("debug message are enabled")

	ssoClient, err := ssogrpc.New(
		context.Background(),
		log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
	)
	if err != nil {
		log.Error("Cannot connect to sso with grpc", sl.Err(err))
		os.Exit(1)
	}
	cls := clients.New(ssoClient)

	dataStorage, err := storage.New(cfg.DataStore)
	if err != nil {
		log.Error("Cannot open storage", sl.Err(err))
		os.Exit(1)
	}

	// Can do with base router net/http
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	{
		router.Post("/", save.New(log, cls, dataStorage))
		router.Delete("/{alias}", delete.New(log, cls, dataStorage))
	}

	router.Get("/{alias}", redirect.New(log, dataStorage))

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
