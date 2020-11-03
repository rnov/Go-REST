package db

import (
	"github.com/rnov/Go-REST/pkg/config"
	"github.com/rnov/Go-REST/pkg/db/redis"
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/rate"
	rcp "github.com/rnov/Go-REST/pkg/recipe"
)

type Recipe interface {
	GetRecipeById(recipeId string) (*rcp.Recipe, error)
	GetAllRecipes() ([]*rcp.Recipe, error)
	CreateRecipe(recipe *rcp.Recipe) error
	UpdateRecipe(recipe *rcp.Recipe) error
	DeleteRecipe(recipeId string) error
}

type Rate interface {
	RateRecipe(recipeId string, rate *rate.Rate) error
}

type Auth interface {
	CheckAuth(auth string) error
}

type Client interface {
	Recipe
	Rate
	Auth
}

func NewDbClient(cfg config.DBConfig) (Client, error) {

	switch cfg.Type {
	case "redis":
		// note: check ping pong etc - consult main
		redisClient := redis.NewRedisClient(cfg.Host, cfg.Port, cfg.Db)
		//check connection with redis
		if _, err := redisClient.Ping().Result(); err != nil {
			return nil, err
		}
		//fmt.Println(pong)
		// create redisProxy with the given client (master)
		return redis.NewRedisProxy(redisClient), nil
	}
	return nil, errors.NewNotFoundErr("")
}
