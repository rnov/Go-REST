package service

import (
	"github.com/rnov/Go-REST/pkg/db"
	r "github.com/rnov/Go-REST/pkg/rate"
	"github.com/rnov/Go-REST/pkg/errors"
)

type Rater interface {
	Rate(id string, rate *r.Rate) error
}

type Service struct {
	rateDb db.Rate
	//// add more func fields - there used to be a logger-
}

func NewRate(rateDb db.Rate) *Service {
	rateSrv := &Service{
		rateDb: rateDb,
	}
	return rateSrv
}

func (r *Service) Rate(id string, rate *r.Rate) error {
	if v := validateRateDataRange(id, rate); len(v) > 0 {
		return errors.NewInvalidParamsErr(v)
	}
	if err := r.rateDb.RateRecipe(id, rate); err != nil {
		return err
	}

	return nil
}

func validateRateDataRange(id string, rate *r.Rate) map[string]string {
	valid := make(map[string]string)

	if len(id) > 100 {
		valid[errors.RateId] = errors.TooLong
	}
	if rate.Note < 1 || rate.Note > 5 {
		valid[errors.Rate] = errors.OutOfRange
	}
	return valid
}
