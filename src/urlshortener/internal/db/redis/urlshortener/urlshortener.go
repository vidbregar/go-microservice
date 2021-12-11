package urlshortener

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/vidbregar/go-microservice/internal/db/models"
)

var ErrFailedSavingUrl = fmt.Errorf("failed saving url")
var ErrFailedLoadingUrl = fmt.Errorf("failed loading url")
var ErrShortPathExists = fmt.Errorf("short path already exist")
var ErrFailedSettingExpire = fmt.Errorf("failed setting expire")

type Storage interface {
	Save(ctx context.Context, shortPath string, item *models.URL) error
	Load(ctx context.Context, shortPath string) (*models.URL, error)
}

type storage struct {
	rdb *redis.Client
}

func New(rdb *redis.Client) Storage {
	return &storage{
		rdb: rdb,
	}
}

func (s *storage) Save(ctx context.Context, shortPath string, url *models.URL) error {
	set, err := s.rdb.HSetNX(ctx, shortPath, "url", url.Url).Result()
	if err != nil {
		return ErrFailedSavingUrl
	}
	if !set {
		return ErrShortPathExists
	}

	err = s.rdb.HSetNX(ctx, shortPath, "expireAt", url.ExpireAt).Err()
	if err != nil {
		return fmt.Errorf("error setting expireAt: %w", ErrFailedSavingUrl)
	}
	if url.ExpireAt > 0 {
		err = s.rdb.ExpireAt(ctx, shortPath, time.Unix(url.ExpireAt, 0)).Err()
		if err != nil {
			return ErrFailedSettingExpire
		}
	}

	return nil
}

func (s *storage) Load(ctx context.Context, shortPath string) (*models.URL, error) {
	res := s.rdb.HGetAll(ctx, shortPath)
	if res.Err() != nil {
		return nil, fmt.Errorf("error loading from Redis: %w", ErrFailedLoadingUrl)
	}

	var item models.URL
	if err := res.Scan(&item); err != nil {
		return nil, fmt.Errorf("error scanning into UrlItem: %w", ErrFailedLoadingUrl)
	}

	return &item, nil
}
