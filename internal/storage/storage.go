package storage

import (
	"errors"
	"url-shortener/internal/config"
	"url-shortener/internal/storage/postgres"
	"url-shortener/internal/storage/redis"

	redisLib "github.com/go-redis/redis"
)

type DataStore struct {
	pgDB      *postgres.Storage
	redisCl   *redis.Cache
	pgAddr    string
	redisAddr string
}

func New(storCnf config.DataStore) (*DataStore, error) {
	pgCfg := storCnf.Storage
	pgDB, err := postgres.New(pgCfg)
	if err != nil {
		return &DataStore{}, err
	}

	redisCfg := storCnf.Cache
	redisCl, err := redis.New(redisCfg)
	if err != nil {
		return &DataStore{}, err
	}

	return &DataStore{
		pgDB:      pgDB,
		redisCl:   redisCl,
		pgAddr:    pgAddr(pgCfg),
		redisAddr: redisCfg.Addr,
	}, nil
}

func (ds *DataStore) GetURL(alias string) (string, error) {
	url, err := ds.redisCl.GetURL(alias)
	if err != nil {
		if errors.Is(err, redisLib.Nil) {
			url, err = ds.pgDB.GetURL(alias)
			if err != nil {
				return "", err
			}
			if err = ds.redisCl.CacheURL(url, alias); err != nil {
				return "", err
			}
			return url, nil
		}
		return "", err
	}
	return url, nil
}

func (ds *DataStore) SaveURL(urlToSave string, alias string) error {
	return ds.pgDB.SaveURL(urlToSave, alias)
}

func (ds *DataStore) DeleteURL(alias string) error {
	return ds.pgDB.DeleteURL(alias)
}

var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists   = errors.New("url exists")
)

func pgAddr(cfg config.Storage) string {
	return cfg.Host + ":" + cfg.Port
}
