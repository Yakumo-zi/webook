package ioc

import "github.com/redis/go-redis/v9"

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
	return cmd
}
