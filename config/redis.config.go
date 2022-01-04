package config

import (
	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
)

func ExampleClient() *redis.Client {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}

func GetDBInstanceRedis() *redis.Client {
	return rdb
}
