package redis

import (
	"fmt"
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/recipes"
)

func (rProxy *RedisProxy) GetRecipeById(recipeId string) (*recipes.Recipe, error) {

	recipeFields, err := rProxy.redisMaster.HGetAll(RedisRecipePattern + recipeId).Result()

	if err != nil {
		return nil, err
	}

	if len(recipeFields) == 0 {
		return nil, errors.NewExistErr(fmt.Sprintf("Does not exist Recipe with Id : %s" + recipeId))
	}

	recipe := mapToRecipeFromRedis(recipeId, recipeFields)
	if recipe == nil {
		return nil, errors.NewExistErr(fmt.Sprint("error parsing int from redis"))
	}

	return recipe, nil

}

func (rProxy *RedisProxy) GetAllRecipes() ([]*recipes.Recipe, error) {

	recipesKeys, err := rProxy.redisMaster.Keys(RedisRecipePattern + RedisAllPattern).Result()

	if err != nil {
		return nil, err
	}

	var recipes []*recipes.Recipe
	for _, key := range recipesKeys {
		redisRcp, err := rProxy.redisMaster.HGetAll(key).Result()
		if err != nil {
			return nil, err
		}

		recipe := mapToRecipeFromRedis(key, redisRcp)
		if recipe == nil {
			return nil, errors.NewDBErr("error parsing int from redis")
			//return nil, errors.New("error parsing int from redis")
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil

}

func (rProxy *RedisProxy) CreateRecipe(recipe *recipes.Recipe) error {

	exists, err := rProxy.redisMaster.Exists(RedisRecipePattern + recipe.Id).Result()
	if err != nil {
		return errors.NewDBErr("")
	}
	if exists > 0 {
		//return errors.New(msg.Exists)
		return errors.NewDBErr("")
	}

	// prepare to insert
	redisFields := mapRecipeToRedisFields(recipe)
	//err = rProxy.HMSet(RedisRecipePattern+recipe.Id, redisFields).Err()
	err = rProxy.redisMaster.HMSet(RedisRecipePattern+recipe.Id, redisFields).Err()

	if err != nil {
		return errors.NewExistErr("")
		//return errors.New(msg.DbError)
	}

	return nil

}

func (rProxy *RedisProxy) UpdateRecipe(recipe *recipes.Recipe) error {

	exists, err := rProxy.redisMaster.Exists(RedisRecipePattern + recipe.Id).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exists == 0 {
		return errors.NewExistErr("")
	}

	// prepare to update
	redisFields := mapRecipeToRedisFields(recipe)
	err = rProxy.redisMaster.HMSet(RedisRecipePattern+recipe.Id, redisFields).Err()

	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	return nil
}

func (rProxy *RedisProxy) DeleteRecipe(recipeId string) error {

	// check whether the recipe has been rated, in that case the rating is also deleted
	exists, err := rProxy.redisMaster.Exists(RedisRecipeRatePattern + recipeId).Result()
	if err != nil {
		//return errors.New(msg.DbError)
		return errors.NewDBErr(err.Error())
	}

	// note todo implement with multi - redis trasactions - (golang redis as : TxPipelined )
	result, err := rProxy.redisMaster.Del(RedisRecipePattern + recipeId).Result()
	//result, err := rProxy.Del(RedisRecipePattern + recipeId).Result()
	if err != nil {
		//return errors.New(msg.DbError)
		return errors.NewDBErr(err.Error())
	}

	if result == 0 {
		//return errors.New(msg.NotFound)
		return errors.NewExistErr("")
	}

	if exists == 1 {
		result, err = rProxy.redisMaster.Del(RedisRecipePattern + recipeId).Result()
		if err != nil {
			return errors.NewDBErr(err.Error())
			//return errors.New(msg.DbError)
		}
	}

	return nil
}

func mapToRecipeFromRedis(key string, redisData map[string]string) *recipes.Recipe {

	//prepTime, err := strconv.Atoi(redisData[msg.Preptime])
	//if err != nil {
	//	return nil
	//}
	//difficulty, err := strconv.Atoi(redisData[msg.Difficulty])
	//if err != nil {
	//	return nil
	//}
	//result := &recipes.Recipe{
	//	Id:         strings.TrimPrefix(key, RedisRecipePattern),
	//	Name:       redisData[msg.Name],
	//	PrepTime:   prepTime,
	//	Difficulty: difficulty,
	//	Vegetarian: redisData[msg.Vegetarian] == RedisTrue,
	//}
	result := &recipes.Recipe{}
	return result
}

func mapRecipeToRedisFields(rcp *recipes.Recipe) map[string]interface{} {
	mappedData := make(map[string]interface{})
	//
	//mappedData[Id] = rcp.Id
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
