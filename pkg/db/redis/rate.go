package redis

import (
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/rate"

	"strconv"
	"time"
)

const (
	allPattern  = "*"
	ratePattern = "RATE_"
)

func (rProxy *Proxy) RateRecipe(recipeID string, rate *rate.Rate) error {
	// check whether recipe exist
	exists, err := rProxy.master.Exists(recipePattern + recipeID).Result()
	if err != nil {
		return errors.NewDBErr(err.Error())
	}

	if exists == 0 {
		return errors.NewExistErr(false)
	}

	// prepare to insert
	redisFields := mapRateToRedisFields(rate.Note)
	err = rProxy.master.HMSet(ratePattern+recipeID, redisFields).Err()

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
