package auth

import (
	"testing"
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
		name   string
		authDB authDBMock
		err    error
	}{
		// todo
		{
			name:   "",
			authDB: authDBMock{},
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			auth := NewAuth(&test.authDB)
			if err := auth.Validate("dXNlcm5hbWU6cGFzc3dvcmQ="); err != test.err {
				t.Errorf("error validation : %w", err)
			}
		})
	}
}

func Test_Encode(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		output string
	}{
		// todo
		{

		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res := encode(test.input)
			if res != test.output {
				t.Errorf("expected output: %s instead got : %s", test.output, res)
			}
		})
	}
}
