package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"url-shotener/internal/config"
	"url-shotener/internal/lib/logger/setup"
	"url-shotener/internal/lib/logger/sl"
)

var info = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"

func main() {
	cfg := config.Mustload()

	log := setup.SetupLogger(cfg.Env)

	log.Info("starting migration", slog.String("env", cfg.Env))

	var pgInfo = fmt.Sprintf(info, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.DBName, cfg.Storage.SSLMode)

	_, err := sql.Open("postgres", pgInfo)
	if err != nil {
		log.Error("error", sl.Err(err))
	}

	if err != nil {
		log.Error("error", sl.Err(err))
	}

	if err != nil {
		log.Error("error", sl.Err(err))
	}
}
