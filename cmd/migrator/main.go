package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"url-shortener/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

func main() {

	var command, migrationDir string
	flag.StringVar(&command, "command", "", "command for migration")
	flag.StringVar(&migrationDir, "dir", "", "migrations dir")

	storage := config.MustLoad().DataStore.Storage

	var dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		storage.Host, storage.Port, storage.User, storage.Password, storage.DBName, storage.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %s", err.Error())
	}
	log.Println("db successfully opened")

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %s", err.Error())
	}
	log.Println("db successfully pinged")

	switch command {
	case "up":
		if err := goose.Up(db, migrationDir); err != nil {
			log.Fatalf("failed to up migrations: %s", err.Error())
		}
		log.Println("migration successfully completed")
	case "down":
		if err := goose.Down(db, migrationDir); err != nil {
			log.Fatalf("failed to down migrations: %s", err.Error())
		}
		log.Println("migration successfully completed")
	case "refresh":
		if err := goose.Down(db, migrationDir); err != nil {
			log.Fatalf("failed to down migrations: %s", err.Error())
		}
		if err := goose.Up(db, migrationDir); err != nil {
			log.Fatalf("failed to up migrations: %s", err.Error())
		}
		log.Println("migration successfully completed")
	default:
		log.Printf("no flag: %s", command)
	}
}
