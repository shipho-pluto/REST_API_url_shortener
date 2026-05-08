package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"url-shotener/internal/config"
	"url-shotener/internal/lib/logger/setup"
	"url-shotener/internal/lib/logger/sl"

	"errors"
	"flag"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var info = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"

func main() {
	var migrationsPath, migrationsTable string

	cfg := config.Mustload()

	log := setup.SetupLogger(cfg.Env)

	log.Info("starting migration", slog.String("env", cfg.Env))

	flag.StringVar(&cfg.Storage.Host, "storage-host", "", "storage host")
	flag.StringVar(&cfg.Storage.Port, "storage-port", "", "storage port")
	flag.StringVar(&cfg.Storage.User, "storage-user", "", "storage user")
	flag.StringVar(&cfg.Storage.Password, "storage-password", "", "storage pass")
	flag.StringVar(&cfg.Storage.DBName, "storage-name", "", "storage name")
	flag.StringVar(&cfg.Storage.SSLMode, "storage-sslmode", "", "storage sslmode")

	var pgInfo = fmt.Sprintf(info, cfg.Storage.Host, cfg.Storage.Port, cfg.Storage.User, cfg.Storage.Password, cfg.Storage.DBName, cfg.Storage.SSLMode)

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	_, err := sql.Open("postgres", pgInfo)
	if err != nil {
		log.Error("error with open database", sl.Err(err))
	}

	if migrationsPath == "" {
		log.Error("migrations-path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("postgres://x-migrations-table=%s", migrationsTable),
	)
	if err != nil {
		log.Error("error", sl.Err(err))
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Error("no migrations to apply")
		}

		log.Error("error", sl.Err(err))
	}

	fmt.Println("migrations applied")
}

// Log represents the logger
type Log struct {
	verbose bool
}

// Printf prints out formatted string into a log
func (l *Log) Printf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

// Verbose shows if verbose print enabled
func (l *Log) Verbose() bool {
	return false
}
