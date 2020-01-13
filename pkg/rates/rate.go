package rates

// note: rename rates and recipes as a single package data ?

import (
	"errors"
	"github.com/goRest/pkg/db"
	e "github.com/goRest/pkg/errors"
	log "github.com/goRest/pkg/logger"
)

type Rate struct {
	Note int `json:"note"`
}

type RateService struct {
	rateDb db.RateDbCalls
	// add more func fields
	logger log.Loggers
}

func NewRateSrv(rateDb db.RateDbCalls, logger log.Loggers) *RateService {
	rateSrv := &RateService{
		rateDb: rateDb,
		logger: logger,
	}
	return rateSrv
}

// fixme errors and pkg/errors should no be use in the same package... erros should be inside pgk/err
func (rateDb *RateService) Rate(id string, rate *Rate) (map[string]string, error) {

	valid := validateRateDataRange(id, rate)
	if len(valid) > 0 {
		return valid, &e.InvalidParamsErr{}
	}

	err := rateDb.rateDb.RateRecipe(id, rate)
	if err != nil {
		if errors.Is(err, &e.DBErr{}) {
			rateDb.logger.Error(err)
		}
		return nil, err
	}

	return nil, nil
}

// fixme these error msg should no be in pkg/errors rather defined in these package
func validateRateDataRange(id string, rate *Rate) map[string]string {
	valid := make(map[string]string)

	if len(id) > 100 {
		valid[e.RateId] = e.TooLong
	}
	if rate.Note < 1 || rate.Note > 5 {
		valid[e.Rate] = e.OutOfRange
	}
	return valid
}
