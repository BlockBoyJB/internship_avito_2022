package redis

import "github.com/redis/go-redis/v9"

type Option func(opts *redis.Options)

func MaxPoolSize(size int) Option {
	return func(redis *redis.Options) {
		redis.PoolSize = size
	}
}

func SetPassword(password string) Option {
	return func(redis *redis.Options) {
		redis.Password = password
	}
}
