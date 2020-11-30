package redis

import (
	e "errors"
	"testing"

	"github.com/rnov/Go-REST/pkg/errors"
)

func TestProxy_CheckAuth(t *testing.T) {
	tests := []struct {
		name        string
		Auth        string
		accessor    *redisAccessorMock
		expectedErr error
	}{
		{
			name: "successful auth",
			Auth: "qwertyzxcv12345",
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 1, nil
				},
			},
		},
		{
			name: "error - exist check from DB",
			Auth: "qwertyzxcv12345",
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, e.New("DB error")
				},
			},
			expectedErr: errors.NewDBErr("DB error"),
		},
		{
			name: "error - auth does not exist",
			Auth: "qwertyzxcv12345",
			accessor: &redisAccessorMock{
				existsAccessor: func(key string) (int64, error) {
					return 0, nil
				},
			},
			expectedErr: errors.NewFailedAuthErr(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			proxy := newRedisMock(test.accessor)
			err := proxy.CheckAuth(test.Auth)
			if err != nil && err.Error() != test.expectedErr.Error() {
				t.Errorf("expected: '%s' instead got: '%s'", test.expectedErr, err)
			}
		})
	}
}
