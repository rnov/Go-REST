package service

import (
	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/errors"
	r "github.com/rnov/Go-REST/pkg/rate"
)

type RateSrv interface {
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
		return errors.NewInputError("Invalid input parameters", v)
	}
	if err := r.rateDb.RateRecipe(id, rate); err != nil {
		return err
	}

	return nil
}

func validateRateDataRange(ID string, rate *r.Rate) map[string]string {
	valid := make(map[string]string)

	if len(ID) > 12 {
		valid[errors.RateID] = errors.TooLong
	}
	if rate.Note < 1 || rate.Note > 5 {
		valid[errors.Rate] = errors.OutOfRange
	}
	return valid
}
