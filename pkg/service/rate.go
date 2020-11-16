package service

import (
	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/errors"
	r "github.com/rnov/Go-REST/pkg/rate"
)

type Rater interface {
	Rate(ID string, rate *r.Rate) error
}

type Rate struct {
	rateDB db.Rate
	//// add more func fields - there used to be a logger-
}

func NewRate(rateDB db.Rate) *Rate {
	rateSrv := &Rate{
		rateDB: rateDB,
	}
	return rateSrv
}

func (r *Rate) Rate(ID string, rate *r.Rate) error {
	if v := validateRateDataRange(ID, rate); len(v) > 0 {
		return errors.NewInputError("Invalid input parameters", v)
	}
	if err := r.rateDB.RateRecipe(ID, rate); err != nil {
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
