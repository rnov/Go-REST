package db

import (
	"errors"

	"github.com/rnov/Go-REST/pkg/config"
	"github.com/rnov/Go-REST/pkg/db/redis"
	"github.com/rnov/Go-REST/pkg/rate"
	rcp "github.com/rnov/Go-REST/pkg/recipe"
)

// Recipe - Provides all DB operations related to recipe's business logic.
type Recipe interface {
	GetRecipeByID(recipeID string) (*rcp.Recipe, error)
	GetAllRecipes() ([]*rcp.Recipe, error)
	CreateRecipe(recipe *rcp.Recipe) error
	UpdateRecipe(recipe *rcp.Recipe) error
	DeleteRecipe(recipeID string) error
}

// Rate - Provides all DB operations related to rate's business logic.
type Rate interface {
	RateRecipe(recipeID string, rate *rate.Rate) error
}

// Auth - Provides all DB operations related to authorization's business logic.
type Auth interface {
	CheckAuth(auth string) error
}

// Client - is a `superset` of DB interfaces that defines a DB client, that way is ensured that a given DB client needs to
//implement all the accessor interfaces.
type Client interface {
	Recipe
	Rate
	Auth
}

// NewClient - DB client constructor based on the configuration that has been loaded.
func NewClient(cfg config.DBConfig) (Client, error) {
	switch cfg.Name {
	case "redis":
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
