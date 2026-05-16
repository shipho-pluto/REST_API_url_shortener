package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"

	"url-shortener/internal/clients"
	"url-shortener/internal/clients/kafka"
	ssogrpc "url-shortener/internal/clients/sso/grpc"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/setup"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/server/handlers/auth/login"
	"url-shortener/internal/server/handlers/auth/register"
	redirect "url-shortener/internal/server/handlers/redirect"
	"url-shortener/internal/server/handlers/url/delete"
	"url-shortener/internal/server/handlers/url/save"
	"url-shortener/internal/server/middleware/auth"
	"url-shortener/internal/server/middleware/logger"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := setup.SetupLogger(cfg.Env)

	log.Info("starting url-shotener service", slog.String("env", cfg.Env))
	log.Debug("debug message are enabled")

	ctx := context.Background()

	ssoClient, err := ssogrpc.New(
		ctx, log,
		cfg.Clients.SSO.Address,
		cfg.Clients.SSO.Timeout,
		cfg.Clients.SSO.RetriesCount,
		cfg.Clients.SSO.AppID,
	)
	if err != nil {
		log.Error("Cannot connect to sso with grpc", sl.Err(err))
		os.Exit(1)
	}

	broker, err := kafka.New(log, cfg.Clients.Broker)
	if err != nil {
		log.Error("Cannot connect to kafka", sl.Err(err))
		os.Exit(1)
	}

	cls := clients.New(ssoClient, broker)

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
	router.Use(auth.New(log))
	router.Use(middleware.Recoverer)

	router.Post("/url/", save.New(log, dataStorage))
	router.Delete("/url/{alias}", delete.New(log, dataStorage))
	router.Post("/url/login/", login.New(ctx, log, cls))
	router.Post("/url/register/", register.New(ctx, log, cls))
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
		if err := cls.Stop(); err != nil {
			log.Error("failed with stop kafka")
		}
	}

	log.Error("server stopped")
}
