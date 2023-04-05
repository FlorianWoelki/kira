package cache

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/go-redis/cache/v8"
)

// Cache is a generic cache implementation that can store and retrieve values of any type.
type Cache[T any] struct {
	// internal represents the internally redis cache.
	internal *cache.Cache
}

// NewCache creates a new cache instance and returns a pointer to it. It also creates a
// new local cache from redis which is a TinyLFU cache.
func NewCache[T any]() *Cache[T] {
	return &Cache[T]{internal: cache.New(&cache.Options{
		Redis:      redisClient,
		LocalCache: cache.NewTinyLFU(1000, time.Minute),
	})}
}

// Set stores a value in the cache for a given language and content key. The TTL is one
// hour. It will also return an error, if an error occured while inserting the value.
func (c *Cache[T]) Set(language, content string, value T) error {
	// Create one key by encoding language and content.
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

// Get retrieves a value from the cache for a given language and content key. It will
// return an error, if an error occured while getting the key or value.
func (c Cache[T]) Get(language, content string) (T, error) {
	// Encode content to get one key.
	key := encodeContent(language + content)

	var wanted T
	if err := c.internal.Get(context.Background(), key, &wanted); err != nil {
		return wanted, err
	}

	return wanted, nil
}

// encodeContent encodes a string into a base64 string.
func encodeContent(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}
