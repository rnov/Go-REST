package auth

import (
	e "errors"
	"testing"

	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/logger"
)

type authDBMock struct {
	checkAuth func(auth string) error
}

func (am *authDBMock) CheckAuth(auth string) error {
	if am.checkAuth != nil {
		return am.checkAuth(auth)
	}
	panic("Not implemented")
}

func TestAuth_Validate(t *testing.T) {
	tests := []struct {
		name        string
		authDB      authDBMock
		expectedErr error
	}{
		{
			name: "successful validation",
			authDB: authDBMock{
				checkAuth: func(auth string) error {
					return nil
				},
			},
			expectedErr: nil,
		},
		{
			name: "error - failed validation",
			authDB: authDBMock{
				checkAuth: func(auth string) error {
					return errors.NewFailedAuthErr()
				},
			},
			expectedErr: errors.NewFailedAuthErr(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := logger.NewLogger()
			auth := NewAuth(&test.authDB, l)
			err := auth.Validate("dXNlcm5hbWU6cGFzc3dvcmQ=")
			if err != nil && !e.Is(err, test.expectedErr) {
				t.Errorf("error validation, unexpected error type")
			}
		})
	}
}
