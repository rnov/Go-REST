package db

import (
	"errors"
	rcp "github.com/goRest/pkg/recipes"
	"github.com/goRest/pkg/rates"
)

const (
	RedisAllPattern        = "*"
	RedisRecipePattern     = "RECIPE_"
	RedisTokenPattern      = "TOKEN_"
	RedisTrue              = "True"
	RedisRecipeRatePattern = "RATE_"
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

func (rProxy *RedisProxy) CheckAuthToken(auth string) error {

	exist, err := rProxy.Exists(RedisTokenPattern + auth).Result()
	if err != nil {
		return err
	}
	if exist == 0 {
		return errors.New(msg.AuthFailed)
	}

	return nil
}
