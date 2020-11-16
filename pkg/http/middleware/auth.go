package middleware

import (
	"net/http"
	"strings"

	"github.com/rnov/Go-REST/pkg/auth"
	"github.com/rnov/Go-REST/pkg/errors"
)

const (
	authHeader = "Authentication"
	basic      = "Basic"
)

func Authentication(auth *auth.Auth, next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
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

func validateAuthStructure(ah string) (string, bool) {
	if res := strings.Split(ah, " "); res[0] != basic && len(res) == 2 {
		return res[1], true
	}
	return "", false
}
