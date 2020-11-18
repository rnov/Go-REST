package redis

import (
	"fmt"

	"github.com/go-redis/redis"
)

// In order to be able to mock redis DB access without 3th parties or running an actual instance.
type redisAccessor interface {
	getAll(key string) (map[string]string, error)
	keys(pattern string) ([]string, error)
	exists(key string) (int64, error)
	set(key string, fields map[string]interface{}) (string, error)
	setErr(key string, fields map[string]interface{}) error
	del(key string) (int64, error)
}

type Proxy struct {
	main redis.Client
	mock redisAccessor
	// could add the workers too e.g: get some data from main and use it to CRUD the workers or register custom logger
}

func NewRedisProxy(client *redis.Client) *Proxy {
	redisProxy := new(Proxy)
	if client != nil {
		redisProxy.main = *client
	}
	return redisProxy
}

func newRedisMock(ra redisAccessor) *Proxy {
	return &Proxy{
		mock: ra,
	}
}

func NewRedisClient(host string, port int, db int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       db,
	})
	return client
}

// to document : a compromise since we have a 3th party and Go "cannot define methods on non-local type" e.g redis.Client is outside the package
func (p *Proxy) getAll(key string) (map[string]string, error) {
	if p.mock != nil {
		return p.mock.getAll(key)
	}
	return p.main.HGetAll(key).Result()
}

func (p *Proxy) keys(pattern string) ([]string, error) {
	if p.mock != nil {
		return p.mock.keys(pattern)
	}
	return p.main.Keys(pattern).Result()
}

func (p *Proxy) exists(key string) (int64, error) {
	if p.mock != nil {
		return p.mock.exists(key)
	}
	return p.main.Exists(key).Result()
}

func (p *Proxy) set(key string, fields map[string]interface{}) (string, error) {
	if p.mock != nil {
		return p.mock.set(key, fields)
	}
	return p.main.HMSet(key, fields).Result()
}

func (p *Proxy) setErr(key string, fields map[string]interface{}) error {
	if p.mock != nil {
		return p.mock.setErr(key, fields)
	}
	return p.main.HMSet(key, fields).Err()
}

func (p *Proxy) del(key string) (int64, error) {
	if p.mock != nil {
		return p.mock.del(key)
	}
	return p.main.Del(key).Result()
}
