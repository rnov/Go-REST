package redis

import (
	"testing"

	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/rate"
)

func TestProxy_RateRecipe(t *testing.T) {
	tests := []struct {
		name        string
		ID          string
		inputRate   *rate.Rate
		accessor    *redisAccessorMock
		expectedErr error
	}{
		{
			name:      "successful rate",
			ID:        "654321",
			inputRate: &rate.Rate{Note: 4},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
				setErrAccessor: func(key string, fields map[string]interface{}) error {
					return nil
				},
			},
		},
		{
			name:      "error - exist check from DB",
			ID:        "654321",
			inputRate: &rate.Rate{Note: 4},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, errors.NewDBErr("DB error")
				},
			},
			expectedErr: errors.NewDBErr("DB error"),
		},
		{
			name:      "error - recipe does not exist",
			ID:        "654321",
			inputRate: &rate.Rate{Note: 4},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, nil
				},
			},
			expectedErr: errors.NewExistErr(false),
		},
		{
			name:      "error - writing rate in DB",
			ID:        "654321",
			inputRate: &rate.Rate{Note: 4},
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
				setErrAccessor: func(key string, fields map[string]interface{}) error {
					return errors.NewDBErr("DB error")
				},
			},
			expectedErr: errors.NewDBErr("DB error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proxy := newRedisMock(test.accessor)
			err := proxy.RateRecipe(test.ID, test.inputRate)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
		})
	}
}
