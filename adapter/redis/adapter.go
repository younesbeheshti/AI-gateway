package redisadapter

import "github.com/redis/go-redis/v9"

type Config struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type RedisAdapter struct {
	client *redis.Client
}

func New() *RedisAdapter {

	//TODO: add to config
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &RedisAdapter{client: client}
}

func (r *RedisAdapter) Client() *redis.Client {
	return r.client
}
