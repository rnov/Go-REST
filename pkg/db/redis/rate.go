package redis

import (
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/rate"

	"strconv"
	"time"
)

const (
	RedisAllPattern        = "*"
	RedisRecipePattern     = "RECIPE_"
	RedisTokenPattern      = "TOKEN_"
	RedisTrue              = "True"
	RedisRecipeRatePattern = "RATE_"
)

func (rProxy *Proxy) RateRecipe(recipeId string, rate *rate.Rate) error {

	// check whether recipe exist
	exists, err := rProxy.master.Exists(RedisRecipePattern + recipeId).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	if exists == 0 {
		return errors.NewNotFoundErr("not found")
	}

	// prepare to insert
	redisFields := mapRateToRedisFields(rate.Note)
	err = rProxy.master.HMSet(RedisRecipeRatePattern+recipeId, redisFields).Err()

	if err != nil {
		return errors.NewDBErr(err.Error())
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

func (rProxy *Proxy) CheckAuthToken(auth string) error {

	exist, err := rProxy.master.Exists(RedisTokenPattern + auth).Result()
	if err != nil {
		return err
	}
	if exist == 0 {
		return errors.NewAuthFailedErr("")
	}

	return nil
}
