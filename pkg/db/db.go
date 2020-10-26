package db

import (
	"github.com/rnov/Go-REST/pkg/db/redis"
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/rate"
	rcp "github.com/rnov/Go-REST/pkg/recipe"
)

type RecipesDbCalls interface {
	GetRecipeById(recipeId string) (*rcp.Recipe, error)
	GetAllRecipes() ([]*rcp.Recipe, error)
	CreateRecipe(recipe *rcp.Recipe) error
	UpdateRecipe(recipe *rcp.Recipe) error
	DeleteRecipe(recipeId string) error
}

type RateDbCalls interface {
	RateRecipe(recipeId string, rate *rate.Rate) error
}

type AuthDb interface {
	CheckAuth(auth string) error
}

type Client interface {
	RecipesDbCalls
	RateDbCalls
	AuthDb
}

func NewDbClient(t string) (Client, error) {

	switch t {
	case "redis":
		// note: check ping pong etc - consult main
		redisClient := redis.NewRedisClient("host", 8080, 1)
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
