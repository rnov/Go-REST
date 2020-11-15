package redis

import (
	"fmt"
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/recipe"
)

const recipePattern = "RECIPE_"

func (rProxy *Proxy) GetRecipeByID(recipeId string) (*recipe.Recipe, error) {
	recipeFields, err := rProxy.master.HGetAll(recipePattern + recipeId).Result()
	if err != nil {
		return nil, err
	}

	if len(recipeFields) == 0 {
		return nil, errors.NewExistErr(false)
	}

	rcp := mapToRecipeFromRedis(recipeId, recipeFields)
	if rcp == nil {
		return nil, errors.NewDBErr(fmt.Sprint("error parsing recipe from redis"))
	}

	return rcp, nil
}

func (rProxy *Proxy) GetAllRecipes() ([]*recipe.Recipe, error) {

	recipesKeys, err := rProxy.master.Keys(recipePattern + allPattern).Result()

	if err != nil {
		return nil, err
	}

	var recipes []*recipe.Recipe
	for _, key := range recipesKeys {
		redisRcp, err := rProxy.master.HGetAll(key).Result()
		if err != nil {
			return nil, err
		}

		recipe := mapToRecipeFromRedis(key, redisRcp)
		if recipe == nil {
			return nil, errors.NewDBErr("error parsing rcp from redis")
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil

}

func (rProxy *Proxy) CreateRecipe(recipe *recipe.Recipe) error {

	exists, err := rProxy.master.Exists(recipePattern + recipe.ID).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exists > 0 {
		return errors.NewExistErr(true)
	}

	// prepare to insert
	redisFields := mapRecipeToRedisFields(recipe)
	_, err = rProxy.master.HMSet(recipePattern+recipe.ID, redisFields).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	return nil

}

func (rProxy *Proxy) UpdateRecipe(recipe *recipe.Recipe) error {

	exists, err := rProxy.master.Exists(recipePattern + recipe.ID).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exists == 0 {
		return errors.NewExistErr(false)
	}

	// prepare to update
	redisFields := mapRecipeToRedisFields(recipe)
	err = rProxy.master.HMSet(recipePattern+recipe.ID, redisFields).Err()

	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	return nil
}

// fixme in case a recipe does not exist: return an existError (new to be created)
func (rProxy *Proxy) DeleteRecipe(recipeId string) error {

	// check whether the recipe has been rated, in that case the rating is also deleted
	exists, err := rProxy.master.Exists(ratePattern + recipeId).Result()
	if err != nil {
		//return errors.New(msg.DbError)
		return errors.NewDBErr(err.Error())
	}

	// note todo implement with multi - redis trasactions - (golang redis as : TxPipelined )
	result, err := rProxy.master.Del(recipePattern + recipeId).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	if result == 0 {
		return errors.NewExistErr(false)
	}

	if exists == 1 {
		result, err = rProxy.master.Del(recipePattern + recipeId).Result()
		if err != nil {
			return errors.NewDBErr(err.Error())
			//return errors.New(msg.DbError)
		}
	}

	return nil
}

func mapToRecipeFromRedis(key string, redisData map[string]string) *recipe.Recipe {

	//prepTime, err := strconv.Atoi(redisData[msg.Preptime])
	//if err != nil {
	//	return nil
	//}
	//difficulty, err := strconv.Atoi(redisData[msg.Difficulty])
	//if err != nil {
	//	return nil
	//}
	//result := &recipes.Recipe{
	//	ID:         strings.TrimPrefix(key, recipePattern),
	//	Name:       redisData[msg.Name],
	//	PrepTime:   prepTime,
	//	Difficulty: difficulty,
	//	Vegetarian: redisData[msg.Vegetarian] == RedisTrue,
	//}
	result := &recipe.Recipe{}
	return result
}

func mapRecipeToRedisFields(rcp *recipe.Recipe) map[string]interface{} {
	mappedData := make(map[string]interface{})
	//
	//mappedData[ID] = rcp.ID
	//mappedData[Preptime] = strconv.Itoa(rcp.PrepTime)
	//mappedData[Difficulty] = strconv.Itoa(rcp.Difficulty)
	//
	//mappedData[Name] = rcp.Name
	//
	//if rcp.Vegetarian {
	//	mappedData[Vegetarian] = "True"
	//} else {
	//	mappedData[Vegetarian] = "False"
	//}

	return mappedData
}
