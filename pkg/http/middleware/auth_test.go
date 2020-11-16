package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rnov/Go-REST/pkg/auth"
	"github.com/rnov/Go-REST/pkg/errors"
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

func TestAuthentication(t *testing.T) {
	tests := []struct {
		name           string
		authDB         authDBMock
		AuthHeader     bool
		Auth           string
		next           func(w http.ResponseWriter, r *http.Request)
		expectedStatus int
	}{
		{
			name: "successful validation",
			authDB: authDBMock{
				checkAuth: func(auth string) error {
					return nil
				},
			},
			AuthHeader: true,
			Auth:       "basic dXNlcm5hbWU6cGFzc3dvcmQ=",
			next: func(w http.ResponseWriter, r *http.Request) {
			},
			expectedStatus: 200,
		},
		{
			name:           "error - not valid auth structure",
			AuthHeader:     true,
			Auth:           "notValidAuth123",
			expectedStatus: 401,
		},
		{
			name:           "error - missing auth header",
			expectedStatus: 401,
		},
		{
			name: "error - non existent auth",
			authDB: authDBMock{
				checkAuth: func(auth string) error {
					return errors.NewFailedAuthErr()
				},
			},
			AuthHeader:     true,
			Auth:           "dXNlcm5hbWU6cGFzc3dvcmQ=",
			expectedStatus: 401,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			a := auth.NewAuth(&test.authDB)

			req, err := http.NewRequest("GET", "/auth", nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Add(authHeader, test.Auth)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/auth", Authentication(a, test.next)).Methods("GET")

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			servicesRouter.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.expectedStatus {
				t.Errorf("handler returned wrong status code: expected %v got %v", test.expectedStatus, rr.Code)
			}
		})
	}
}
