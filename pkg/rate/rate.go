package rate

import (
	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/errors"
	log "github.com/rnov/Go-REST/pkg/logger"
)

type Rate struct {
	Note int `json:"note"`
}

type Rater interface {
	Rate(id string, rate *Rate) error
}

type Service struct {
	rateDb db.RateDbCalls
	// add more func fields
	logger log.Loggers
}

func NewRateSrv(rateDb db.RateDbCalls, logger log.Loggers) *Service {
	rateSrv := &Service{
		rateDb: rateDb,
		logger: logger,
	}
	return rateSrv
}

func (rateDb *Service) Rate(id string, rate *Rate) error {
	if v := validateRateDataRange(id, rate); len(v) > 0 {
		return errors.NewInvalidParamsErr(v)
	}
	if err := rateDb.rateDb.RateRecipe(id, rate); err != nil {
		return err
	}

	return nil
}

func validateRateDataRange(id string, rate *Rate) map[string]string {
	valid := make(map[string]string)

	if len(id) > 100 {
		valid[errors.RateId] = errors.TooLong
	}
	if rate.Note < 1 || rate.Note > 5 {
		valid[errors.Rate] = errors.OutOfRange
	}
	return valid
}
