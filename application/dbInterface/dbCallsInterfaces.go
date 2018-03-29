package dbInterface

import (
	msg "Go-REST/application/common"
	"Go-REST/application/model"
	"errors"
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
	GetRecipeById(recipeId string) (*model.Recipe, error)
	GetAllRecipes() ([]*model.Recipe, error)
	CreateRecipe(recipe *model.Recipe) error
	UpdateRecipe(recipe *model.Recipe) error
	DeleteRecipe(recipeId string) error
}

type RateDbCalls interface {
	CheckAuthToken(auth string) error
	RateRecipe(recipeId string, rate *model.Rate) error
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
