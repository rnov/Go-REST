package redis

import (
	"strconv"
	"strings"

	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/recipe"
)

const recipePattern = "RECIPE_"

func (p *Proxy) GetRecipeByID(ID string) (*recipe.Recipe, error) {
	//recipeFields, err := p.main.HGetAll(recipePattern + Auth).Result()
	recipeFields, err := p.getAll(recipePattern + ID)
	if err != nil {
		return nil, err
	}
	if len(recipeFields) == 0 {
		return nil, errors.NewExistErr(false)
	}

	rcp := mapToRecipeFromRedis(ID, recipeFields)
	if rcp == nil {
		return nil, errors.NewDBErr("error parsing recipe from redis")
	}

	return rcp, nil
}

func (p *Proxy) GetAllRecipes() ([]*recipe.Recipe, error) {
	//recipesKeys, err := p.main.Keys(recipePattern + allPattern).Result()
	recipesKeys, err := p.keys(recipePattern + allPattern)
	if err != nil {
		return nil, err
	}

	recipes := make([]*recipe.Recipe, 0)
	for _, key := range recipesKeys {
		//redisRcp, err := p.main.HGetAll(key).Result()
		redisRcp, err := p.getAll(key)
		if err != nil {
			return nil, err
		}

		rcp := mapToRecipeFromRedis(key, redisRcp)
		if rcp == nil {
			return nil, errors.NewDBErr("error parsing rcp from redis")
		}
		recipes = append(recipes, rcp)
	}

	return recipes, nil
}

func (p *Proxy) CreateRecipe(recipe *recipe.Recipe) error {
	//exists, err := p.main.Exists(recipePattern + recipe.Auth).Result()
	exists, err := p.exists(recipePattern + recipe.ID)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exists > 0 {
		return errors.NewExistErr(true)
	}

	// prepare to insert
	redisFields := mapRecipeToRedisFields(recipe)
	//_, err = p.main.HMSet(recipePattern+recipe.Auth, redisFields).Result()
	_, err = p.set(recipePattern+recipe.ID, redisFields)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	return nil
}

func (p *Proxy) UpdateRecipe(recipe *recipe.Recipe) error {
	//exists, err := p.main.Exists(recipePattern + recipe.Auth).Result()
	exists, err := p.exists(recipePattern + recipe.ID)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exists == 0 {
		return errors.NewExistErr(false)
	}

	// prepare to update
	redisFields := mapRecipeToRedisFields(recipe)
	//err = p.main.HMSet(recipePattern+recipe.ID, redisFields).Err()
	err = p.setErr(recipePattern+recipe.ID, redisFields)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	return nil
}

func (p *Proxy) DeleteRecipe(ID string) error {
	// check whether the recipe has been rated, in that case the rating is also deleted
	//exists, err := p.main.Exists(ratePattern + Auth).Result()
	exists, err := p.exists(ratePattern + ID)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	// note todo implement with multi - redis trasactions - (golang redis as : TxPipelined )
	//result, err := p.main.Del(recipePattern + Auth).Result()
	result, err := p.del(recipePattern + ID)
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	if result == 0 {
		return errors.NewExistErr(false)
	}

	if exists == 1 {
		//_, err = p.main.Del(recipePattern + Auth).Result()
		_, err = p.del(recipePattern + ID)
		if err != nil {
			return errors.NewDBErr(err.Error())
		}
	}

	return nil
}

func mapToRecipeFromRedis(key string, redisData map[string]string) *recipe.Recipe {
	prepTime, err := strconv.Atoi(redisData["preptime"])
	if err != nil {
		return nil
	}
	difficulty, err := strconv.Atoi(redisData["difficulty"])
	if err != nil {
		return nil
	}
	result := &recipe.Recipe{
		ID:         strings.TrimPrefix(key, recipePattern),
		Name:       redisData["Name"],
		PrepTime:   prepTime,
		Difficulty: difficulty,
		Vegetarian: redisData["Vegetarian"] == "TRUE",
	}
	return result
}

func mapRecipeToRedisFields(rcp *recipe.Recipe) map[string]interface{} {
	mappedData := make(map[string]interface{})

	mappedData["ID"] = rcp.ID
	mappedData["Preptime"] = strconv.Itoa(rcp.PrepTime)
	mappedData["Difficulty"] = strconv.Itoa(rcp.Difficulty)

	mappedData["Name"] = rcp.Name

	if rcp.Vegetarian {
		mappedData["Vegetarian"] = "True"
	} else {
		mappedData["Vegetarian"] = "False"
	}

	return mappedData
}
