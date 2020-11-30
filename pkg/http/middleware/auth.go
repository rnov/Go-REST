package middleware

import (
	"net/http"
	"strings"

	"github.com/rnov/Go-REST/pkg/auth"
	"github.com/rnov/Go-REST/pkg/errors"
)

const (
	authHeader = "Authorization"
	basic      = "Basic"
)

// Authentication - custom HTTP middleware that validates user's basic auth.
func Authentication(auth auth.Validator, next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var basicAuth string
		ah := r.Header.Get(authHeader)
		basicAuth, valid := validateAuthStructure(ah)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if err := auth.Validate(basicAuth); err != nil {
			errors.BuildResponse(w, r.Method, err)
			return
		}
		next(w, r)
	}
}

// validateAuthStructure - validates that the authorization header value provided by the user has a valid structure.
func validateAuthStructure(ah string) (string, bool) {
	if res := strings.Split(ah, " "); res[0] == basic && len(res) == 2 {
		return res[1], true
	}

	return "", false
}
