package middleware

import (
	"github.com/rnov/Go-REST/pkg/auth"
	"github.com/rnov/Go-REST/pkg/http/rest"
	"net/http"
)

const (
	authHeader = "Authentication"
	basic      = "Basic"
)

func Authentication(auth *auth.Auth, next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ah := r.Header[authHeader]
		if len(ah) != 1 && !validateAuthStructure(ah) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		basicAuth := ah[1]
		if err := auth.Validate(basicAuth); err != nil {
			rest.BuildErrorResponse(w, err)
			return
		}
		next(w, r)
	}
}

func validateAuthStructure(ah []string) bool {
	return len(ah) != 2 || ah[0] != basic
}
