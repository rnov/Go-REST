package redis

import (
	"fmt"
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/recipe"
)

func (rProxy *Proxy) GetRecipeById(recipeId string) (*recipe.Recipe, error) {

	recipeFields, err := rProxy.master.HGetAll(RedisRecipePattern + recipeId).Result()

	if err != nil {
		return nil, err
	}

	if len(recipeFields) == 0 {
		return nil, errors.NewExistErr(fmt.Sprintf("Does not exist Recipe with ID : %s" + recipeId))
	}

	recipe := mapToRecipeFromRedis(recipeId, recipeFields)
	if recipe == nil {
		return nil, errors.NewExistErr(fmt.Sprint("error parsing int from redis"))
	}

	return recipe, nil

}

func (rProxy *Proxy) GetAllRecipes() ([]*recipe.Recipe, error) {

	recipesKeys, err := rProxy.master.Keys(RedisRecipePattern + RedisAllPattern).Result()

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
			return nil, errors.NewDBErr("error parsing int from redis")
			//return nil, errors.New("error parsing int from redis")
		}
		recipes = append(recipes, recipe)
	}

	return recipes, nil

}

func (rProxy *Proxy) CreateRecipe(recipe *recipe.Recipe) error {

	exists, err := rProxy.master.Exists(RedisRecipePattern + recipe.ID).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exists > 0 {
		return errors.NewExistErr(fmt.Sprintf("recipe with ID %s already exists", recipe.ID))
	}

	// prepare to insert
	redisFields := mapRecipeToRedisFields(recipe)
	_, err = rProxy.master.HMSet(RedisRecipePattern+recipe.ID, redisFields).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	return nil

}

func (rProxy *Proxy) UpdateRecipe(recipe *recipe.Recipe) error {

	exists, err := rProxy.master.Exists(RedisRecipePattern + recipe.ID).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}
	if exists == 0 {
		return errors.NewExistErr("")
	}

	// prepare to update
	redisFields := mapRecipeToRedisFields(recipe)
	err = rProxy.master.HMSet(RedisRecipePattern+recipe.ID, redisFields).Err()

	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	return nil
}

func (rProxy *Proxy) DeleteRecipe(recipeId string) error {

	// check whether the recipe has been rated, in that case the rating is also deleted
	exists, err := rProxy.master.Exists(RedisRecipeRatePattern + recipeId).Result()
	if err != nil {
		//return errors.New(msg.DbError)
		return errors.NewDBErr(err.Error())
	}

	// note todo implement with multi - redis trasactions - (golang redis as : TxPipelined )
	result, err := rProxy.master.Del(RedisRecipePattern + recipeId).Result()
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
		result, err = rProxy.master.Del(RedisRecipePattern + recipeId).Result()
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
	//	ID:         strings.TrimPrefix(key, RedisRecipePattern),
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
