package db

import (
	"github.com/rnov/Go-REST/pkg/db/redis"
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/rates"
	rcp "github.com/rnov/Go-REST/pkg/recipes"
)

type RecipesDbCalls interface {
	CheckAuthToken(auth string) error
	GetRecipeById(recipeId string) (*rcp.Recipe, error)
	GetAllRecipes() ([]*rcp.Recipe, error)
	CreateRecipe(recipe *rcp.Recipe) error
	UpdateRecipe(recipe *rcp.Recipe) error
	DeleteRecipe(recipeId string) error
}

type RateDbCalls interface {
	CheckAuthToken(auth string) error
	RateRecipe(recipeId string, rate *rates.Rate) error
}

type Client interface {
	RecipesDbCalls
	RateDbCalls
}

func NewDbClient(t string) (Client, error) {
	//var retVal = RecipesDbCalls{}
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
