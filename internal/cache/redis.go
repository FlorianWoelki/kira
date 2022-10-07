package cache

import "github.com/go-redis/redis/v8"

var redisClient = redis.NewClient(&redis.Options{
	Addr:     "cache:6379",
	DB:       0,
	Password: "",
})
