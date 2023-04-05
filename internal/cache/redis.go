package cache

import "github.com/go-redis/redis/v8"

// redisClient is an internally used cache and connected to port `6379`.
var redisClient = redis.NewClient(&redis.Options{
	Addr:     "cache:6379",
	DB:       0,
	Password: "",
})
