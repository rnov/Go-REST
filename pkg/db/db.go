package db

import (
	"errors"
	"github.com/rnov/Go-REST/pkg/config"
	"github.com/rnov/Go-REST/pkg/db/redis"
	"github.com/rnov/Go-REST/pkg/rate"
	rcp "github.com/rnov/Go-REST/pkg/recipe"
)

type Recipe interface {
	GetRecipeByID(recipeID string) (*rcp.Recipe, error)
	GetAllRecipes() ([]*rcp.Recipe, error)
	CreateRecipe(recipe *rcp.Recipe) error
	UpdateRecipe(recipe *rcp.Recipe) error
	DeleteRecipe(recipeID string) error
}

type Rate interface {
	RateRecipe(recipeID string, rate *rate.Rate) error
}

type Auth interface {
	CheckAuth(auth string) error
}

type Client interface {
	Recipe
	Rate
	Auth
}

func NewClient(cfg config.DBConfig) (Client, error) {
	switch cfg.Name {
	case "redis":
		// note: check ping pong etc - consult main
		redisClient := redis.NewRedisClient(cfg.Host, cfg.Port, cfg.DB)
		//check connection with redis
		if _, err := redisClient.Ping().Result(); err != nil {
			return nil, err
		}
		//fmt.Println(pong)
		return redis.NewRedisProxy(redisClient), nil
	}
	return nil, errors.New("database does not exist")
}
