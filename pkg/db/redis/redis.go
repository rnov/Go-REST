package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

type Proxy struct {
	master redis.Client
	// could add the slaves too e.g: get some data from master and use it to CRUD the slaves or register custom logger
}

func NewRedisProxy(client *redis.Client) *Proxy {
	redisProxy := new(Proxy)
	redisProxy.master = *client
	return redisProxy
}

func NewRedisClient(host string, port int, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})
	return client
}
