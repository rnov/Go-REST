package dbInterface

import (
	"fmt"

	"github.com/go-redis/redis"
)

type RedisProxy struct {
	redisMaster redis.Client
	// could add the slaves too e.g: get some data from master and use it to CRUD the slaves or register custom logger
}

func NewRedisProxy(client *redis.Client) *RedisProxy {
	redisProxy := new(RedisProxy)
	redisProxy.redisMaster = *client
	return redisProxy
}

func NewClient(host string, port int, db int) *redis.Client {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})
	return client
}

type RedisWrapper interface {
	Del(key string) *redis.IntCmd
	HGetAll(key string) *redis.StringStringMapCmd
	HMSet(key string, fields map[string]interface{}) *redis.StatusCmd
	Exists(key string) *redis.IntCmd
	Keys(pattern string) *redis.StringSliceCmd
}

func (rp *RedisProxy) Del(key string) *redis.IntCmd {
	return rp.redisMaster.Del(key)
}

func (rp *RedisProxy) HGetAll(key string) *redis.StringStringMapCmd {
	return rp.redisMaster.HGetAll(key)
}

// fields map[string]interface{}
func (rp *RedisProxy) HMSet(key string, fields map[string]interface{}) *redis.StatusCmd {
	return rp.redisMaster.HMSet(key, fields)
}

func (rp *RedisProxy) Exists(key string) *redis.IntCmd {
	return rp.redisMaster.Exists(key)
}

func (rp *RedisProxy) Keys(pattern string) *redis.StringSliceCmd {
	return rp.redisMaster.Keys(pattern)
}
