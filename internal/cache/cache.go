package cache

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/go-redis/cache/v8"
)

type Cache[T any] struct {
	internal *cache.Cache
}

func NewCache[T any]() *Cache[T] {
	return &Cache[T]{internal: cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})}
}

func (c *Cache[T]) Set(language, content string, value T) error {
	key := encodeContent(language + content)
	if err := c.internal.Set(&cache.Item{
		Ctx:   context.Background(),
		Key:   key,
		Value: value,
		TTL:   time.Hour,
	}); err != nil {
		return err
	}

	return nil
}

func (c Cache[T]) Get(language, content string) (T, error) {
	key := encodeContent(language + content)

	var wanted T
	if err := c.internal.Get(context.Background(), key, &wanted); err != nil {
		return wanted, err
	}

	return wanted, nil
}

func encodeContent(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}

func decodeHash(hash string) (string, error) {
	s, err := base64.StdEncoding.DecodeString(hash)
	if err != nil {
		return "", err
	}

	return string(s), nil
}
