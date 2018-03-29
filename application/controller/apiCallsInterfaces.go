package controller

import (
	"Go-REST/application/model"
)

// This file is used as interface declarer, here are put all the interfaces that make queries to DBs

type ApiRecipeCalls interface {
	GetById(recipeId string) (*model.Recipe, error)
	ListAll() ([]*model.Recipe, error)
	Create(recipe *model.Recipe, auth string) (map[string]string, error)
	Update(recipe *model.Recipe, urlId string, auth string) (map[string]string, error)
	Delete(recipeId string, auth string) error
}

type ApiRateCalls interface {
	Rate(id string, rating *model.Rate) (map[string]string, error) // rate recipe
}
