package redis

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/recipe"
)

const (
	recipePattern = "RECIPE_"
	rcpID         = "ID"
	name          = "Name"
	prepTime      = "Preptime"
	vegetarian    = "Vegetarian"
	difficulty    = "Difficulty"
)

func (p *Proxy) GetRecipeByID(ID string) (*recipe.Recipe, error) {
	recipeFields, err := p.getAll(recipePattern + ID)
	if err != nil {
		return nil, errors.NewDBErr(err.Error())
	}
	if len(recipeFields) == 0 {
		return nil, errors.NewExistErr(false)
	}
	rcp, err := mapToRecipeFromRedis(ID, recipeFields)
	if err != nil {
		return nil, err
	}

	return rcp, nil
}

func (p *Proxy) GetAllRecipes() ([]*recipe.Recipe, error) {
	recipesKeys, err := p.keys(recipePattern + allPattern)
	if err != nil {
		return nil, errors.NewDBErr(err.Error())
	}

	recipes := make([]*recipe.Recipe, 0)
	for _, key := range recipesKeys {
		redisRcp, err := p.getAll(key)
		if err != nil {
			return nil, errors.NewDBErr(err.Error())
		}
		rcp, err := mapToRecipeFromRedis(key, redisRcp)
		if err != nil {
			return nil, err
		}
		recipes = append(recipes, rcp)
	}

	return recipes, nil
}

func (p *Proxy) CreateRecipe(recipe *recipe.Recipe) error {
	exists, err := p.exists(recipePattern + recipe.ID)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exists > 0 {
		return errors.NewExistErr(true)
	}
	// prepare to insert
	redisFields := mapRecipeToRedisFields(recipe)
	_, err = p.set(recipePattern+recipe.ID, redisFields)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	return nil
}

func (p *Proxy) UpdateRecipe(recipe *recipe.Recipe) error {
	exists, err := p.exists(recipePattern + recipe.ID)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exists == 0 {
		return errors.NewExistErr(false)
	}
	// prepare to update
	redisFields := mapRecipeToRedisFields(recipe)
	err = p.setErr(recipePattern+recipe.ID, redisFields)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	return nil
}

func (p *Proxy) DeleteRecipe(ID string) error {
	// check whether the recipe has been rated, in that case the rating is also deleted
	exists, err := p.exists(ratePattern + ID)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	result, err := p.del(recipePattern + ID)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if result == 0 {
		return errors.NewExistErr(false)
	}
	if exists == 1 {
		_, err = p.del(recipePattern + ID)
		if err != nil {
			return errors.NewDBErr(err.Error())
		}
	}

	return nil
}

func mapToRecipeFromRedis(key string, redisData map[string]string) (*recipe.Recipe, error) {
	prepTime, err := strconv.Atoi(redisData[prepTime])
	if err != nil {
		return nil, errors.NewDBErr(fmt.Sprintf("error parsing recipe from redis: %s", err.Error()))
	}
	difficulty, err := strconv.Atoi(redisData[difficulty])
	if err != nil {
		return nil, errors.NewDBErr(fmt.Sprintf("error parsing recipe from redis: %s", err.Error()))
	}
	result := &recipe.Recipe{
		ID:         strings.TrimPrefix(key, recipePattern),
		Name:       redisData[name],
		PrepTime:   prepTime,
		Difficulty: difficulty,
		Vegetarian: redisData[vegetarian] == "TRUE",
	}

	return result, nil
}

func mapRecipeToRedisFields(rcp *recipe.Recipe) map[string]interface{} {
	mappedData := make(map[string]interface{})
	mappedData[rcpID] = rcp.ID
	mappedData[prepTime] = strconv.Itoa(rcp.PrepTime)
	mappedData[difficulty] = strconv.Itoa(rcp.Difficulty)
	mappedData[name] = rcp.Name
	if rcp.Vegetarian {
		mappedData[vegetarian] = "True"
	} else {
		mappedData[vegetarian] = "False"
	}

	return mappedData
}
