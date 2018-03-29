package controller

import (
	"Go-REST/application/common"
	"Go-REST/application/dbInterface"
	"Go-REST/application/model"
	"errors"
)

type RateController struct {
	rateDb dbInterface.RateDbCalls
	logger common.LogInterface
	// add more func fields
}

func NewRateController(rateDb dbInterface.RateDbCalls, logger common.LogInterface) *RateController {
	rateController := &RateController{
		rateDb: rateDb,
		logger: logger,
	}
	return rateController
}

func (rateDb *RateController) Rate(id string, rate *model.Rate) (map[string]string, error) {

	valid := validateRateDataRange(id, rate)
	if len(valid) > 0 {
		return valid, errors.New(common.InvalidParams)
	}

	err := rateDb.rateDb.RateRecipe(id, rate)
	if err != nil {
		if err.Error() == common.DbError {
			rateDb.logger.Error(err)
		}
		return nil, err
	}

	return nil, nil
}

func validateRateDataRange(id string, rate *model.Rate) map[string]string {
	valid := make(map[string]string)

	if len(id) > 100 {
		valid[common.RateId] = common.TooLong
	}
	if rate.Note < 1 || rate.Note > 5 {
		valid[common.Rate] = common.OutOfRange
	}
	return valid
}
