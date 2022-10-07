package cache

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/go-redis/cache/v8"
)

type Cache struct {
	internal *cache.Cache
}

func NewCache() *Cache {
	return &Cache{internal: cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})}
}

func (c *Cache) Set(key, content string) error {
	value := encodeContent(content)
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

func (c Cache) Get(key string) (string, error) {
	var wanted string
	if err := c.internal.Get(context.Background(), key, &wanted); err != nil {
		return "", err
	}

	value, err := decodeHash(wanted)
	if err != nil {
		return "", err
	}

	return value, nil
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
