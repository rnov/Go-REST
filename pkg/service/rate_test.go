package service

import (
	e "errors"
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/rate"
	"strings"
	"testing"
)

type rateDBMock struct {
	rateRecipe func(recipeId string, rate *rate.Rate) error
}

func (rm *rateDBMock) RateRecipe(recipeId string, rate *rate.Rate) error {
	if rm.rateRecipe != nil {
		return rm.rateRecipe(recipeId, rate)
	}
	panic("Not implemented")
}

func TestService_Rate(t *testing.T) {
	tests := []struct {
		name        string
		rateDB      rateDBMock
		inputRate   rate.Rate
		inputID     string
		expectedErr error
	}{
		{
			name: "successful rate",
			rateDB: rateDBMock{
				rateRecipe: func(recipeId string, rate *rate.Rate) error {
					return nil
				},
			},
			inputRate: rate.Rate{
				Note: 5,
			},
			expectedErr: nil,
		},
		{
			name: "error validating rage: rate ID too long",
			inputRate: rate.Rate{
				Note: 5,
			},
			inputID:     "12sw2329cwme9",
			expectedErr: e.New("invalid parameters"),
		},
		{
			name: "error validating range: note out of range",
			inputRate: rate.Rate{
				Note: 100,
			},
			expectedErr: e.New("invalid parameters"),
		},
		{
			name: "error DB",
			rateDB: rateDBMock{
				rateRecipe: func(recipeId string, rate *rate.Rate) error {
					return errors.NewDBErr("DB error")
				},
			},
			inputRate: rate.Rate{
				Note: 5,
			},
			expectedErr: errors.NewDBErr("DB error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rateSvr := NewRate(&test.rateDB)
			if err := rateSvr.Rate(test.inputID, &test.inputRate); err != nil && !strings.Contains(err.Error(), test.expectedErr.Error()) {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
		})
	}
}
