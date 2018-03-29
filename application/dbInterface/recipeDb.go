package dbInterface

import (
	msg "Go-REST/application/common"
	"Go-REST/application/model"
	"errors"
	"strconv"
	"strings"
)

func (rProxy *RedisProxy) GetRecipeById(recipeId string) (*model.Recipe, error) {

	recipeFields, err := rProxy.HGetAll(RedisRecipePattern + recipeId).Result()

	if err != nil {
		return nil, err
	}

	if len(recipeFields) == 0 {
		return nil, errors.New("Does not exist Recipe with Id : " + recipeId)
	}

	recipe := mapToRecipeFromRedis(recipeId, recipeFields)
	if recipe == nil {
		return nil, errors.New("error parsing int from redis")
	}

	return recipe, nil

}

func (rProxy *RedisProxy) GetAllRecipes() ([]*model.Recipe, error) {

	recipesKeys, err := rProxy.Keys(RedisRecipePattern + RedisAllPattern).Result()

	if err != nil {
		return nil, err
	}

	var recipes []*model.Recipe
	for _, key := range recipesKeys {
		redisRcp, err := rProxy.HGetAll(key).Result()
		if err != nil {
			return nil, err
		}

		recipe := mapToRecipeFromRedis(key, redisRcp)
		if recipe == nil {
			return nil, errors.New("error parsing int from redis")
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil

}

func (rProxy *RedisProxy) CreateRecipe(recipe *model.Recipe) error {

	exists, err := rProxy.Exists(RedisRecipePattern + recipe.Id).Result()
	if err != nil {
		return errors.New(msg.DbError)
	}
	if exists > 0 {
		return errors.New(msg.Exists)
	}

	// prepare to insert
	redisFields := mapRecipeToRedisFields(recipe)
	err = rProxy.HMSet(RedisRecipePattern+recipe.Id, redisFields).Err()

	if err != nil {
		return errors.New(msg.DbError)
	}

	return nil

}

func (rProxy *RedisProxy) UpdateRecipe(recipe *model.Recipe) error {

	exists, err := rProxy.Exists(RedisRecipePattern + recipe.Id).Result()
	if err != nil {
		return errors.New(msg.DbError)
	}
	if exists == 0 {
		return errors.New(msg.NotFound)
	}

	// prepare to update
	redisFields := mapRecipeToRedisFields(recipe)
	err = rProxy.HMSet(RedisRecipePattern+recipe.Id, redisFields).Err()

	if err != nil {
		return errors.New(msg.DbError)
	}

	return nil
}

func (rProxy *RedisProxy) DeleteRecipe(recipeId string) error {

	// check whether the recipe has been rated, in that case the rating is also deleted
	exists, err := rProxy.Exists(RedisRecipeRatePattern + recipeId).Result()
	if err != nil {
		return errors.New(msg.DbError)
	}

	// note todo implement with multi - redis trasactions - (golang redis as : TxPipelined )
	result, err := rProxy.Del(RedisRecipePattern + recipeId).Result()
	if err != nil {
		return errors.New(msg.DbError)
	}

	if result == 0 {
		return errors.New(msg.NotFound)
	}

	if exists == 1 {
		result, err = rProxy.Del(RedisRecipeRatePattern + recipeId).Result()
		if err != nil {
			return errors.New(msg.DbError)
		}
	}

	return nil
}

func mapToRecipeFromRedis(key string, redisData map[string]string) *model.Recipe {

	prepTime, err := strconv.Atoi(redisData[msg.Preptime])
	if err != nil {
		return nil
	}
	difficulty, err := strconv.Atoi(redisData[msg.Difficulty])
	if err != nil {
		return nil
	}
	result := &model.Recipe{
		Id:         strings.TrimPrefix(key, RedisRecipePattern),
		Name:       redisData[msg.Name],
		PrepTime:   prepTime,
		Difficulty: difficulty,
		Vegetarian: redisData[msg.Vegetarian] == RedisTrue,
	}
	return result
}

func mapRecipeToRedisFields(rcp *model.Recipe) map[string]interface{} {
	mappedData := make(map[string]interface{})

	mappedData[msg.Id] = rcp.Id
	mappedData[msg.Preptime] = strconv.Itoa(rcp.PrepTime)
	mappedData[msg.Difficulty] = strconv.Itoa(rcp.Difficulty)

	mappedData[msg.Name] = rcp.Name

	if rcp.Vegetarian {
		mappedData[msg.Vegetarian] = "True"
	} else {
		mappedData[msg.Vegetarian] = "False"
	}

	return mappedData
}
