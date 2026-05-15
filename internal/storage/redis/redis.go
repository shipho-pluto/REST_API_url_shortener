package redis

import (
	"errors"
	"fmt"
	"time"
	"url-shortener/internal/config"

	"github.com/go-redis/redis"
)

type Cache struct {
	cl *redis.Client
}

var expirationTime = time.Minute

func New(cfg config.Cache) (*Cache, error) {
	const op = "storage.redis.Init"

	cl := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		MaxRetries:   cfg.MaxRetries,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
	})

	if err := cl.Ping().Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Cache{
		cl: cl,
	}, nil
}

func (r *Cache) CacheURL(urlToSave string, alias string) error {
	const op = "storage.redis.SaveURL"

	err := r.cl.Set(alias, urlToSave, expirationTime).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Cache) GetURL(alias string) (string, error) {
	const op = "storage.redis.GetURL"

	var url string

	err := r.cl.Get(alias).Scan(&url)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", redis.Nil
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url, nil
}
