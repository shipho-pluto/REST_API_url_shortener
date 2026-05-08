package postgres

import (
	"database/sql"
	"fmt"
	"url-shortener/internal/config"
	req "url-shortener/internal/storage/postgres/requests"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

var info = "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s"

func New(storage *config.Storage) (*Storage, error) {
	const op = "storage.postgres.Init"

	var pgInfo = fmt.Sprintf(info, storage.Host, storage.Port, storage.User, storage.Password, storage.DBName, storage.SSLMode)

	db, err := sql.Open("postgres", pgInfo)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	//TODO: rewrite to migrator

	// stmt, err := db.Prepare(`
	// CREATE TABLE IF NOT EXISTS url (
	// 	id SERIAL PRIMARY KEY,
	// 	alias TEXT UNIQUE,
	// 	url TEXT);
	// `)

	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	// if _, err = stmt.Exec(); err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	// stmt, err = db.Prepare(`
	// CREATE INDEX IF NOT EXISTS idx_alias ON url (alias);
	// `)

	// if err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	// if _, err = stmt.Exec(); err != nil {
	// 	return nil, fmt.Errorf("%s: %w", op, err)
	// }

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave string, alias string) error {
	const op = "storage.postgres.SaveURL"

	stmt, err := s.db.Prepare(req.SaveURLReq)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(urlToSave, alias); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) DeleteURL(alias string) error {
	const op = "storage.postgres.DeleteURL"

	stmt, err := s.db.Prepare(req.DeleteURLReq)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if _, err = stmt.Exec(alias); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const op = "storage.postgres.GetURL"

	stmt, err := s.db.Prepare(req.GetURLReq)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	var url string

	if err = stmt.QueryRow(alias).Scan(&url); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
