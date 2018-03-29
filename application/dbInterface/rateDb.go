package dbInterface

import (
	msg "Go-REST/application/common"
	"Go-REST/application/model"
	"errors"
	"strconv"
	"time"
)

func (rProxy *RedisProxy) RateRecipe(recipeId string, rate *model.Rate) error {

	// check whether recipe exist
	exists, err := rProxy.Exists(RedisRecipePattern + recipeId).Result()
	if err != nil {
		return errors.New(msg.DbError)
	}

	if exists == 0 {
		return errors.New(msg.NotFound)
	}

	// prepare to insert
	redisFields := mapRateToRedisFields(rate.Note)
	err = rProxy.HMSet(RedisRecipeRatePattern+recipeId, redisFields).Err()

	if err != nil {
		return errors.New(msg.DbError)
	}

	return nil
}

func mapRateToRedisFields(rating int) map[string]interface{} {
	mappedData := make(map[string]interface{})
	// since AUTH it is not necessary we use the timestamp as key to insert the rating into redis
	key := strconv.FormatInt(time.Now().Unix(), 10)
	mappedData[key] = rating

	return mappedData
}
