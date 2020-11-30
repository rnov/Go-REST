package auth

import (
	"crypto/sha256"
	"fmt"

	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/logger"
)

// Auth is business logic struct for the authorization custom middleware/
type Auth struct {
	DB  db.Auth
	Log logger.Loggers
}

func NewAuth(db db.Auth, l logger.Loggers) *Auth {
	return &Auth{
		DB:  db,
		Log: l,
	}
}

// Validator - defines all the business logic operations for authorization.
type Validator interface {
	Validate(ba string) error
}

// Validate - validates that a given user is authorized to perform the request, hashes and compares the result with the ones
// stored.
func (a *Auth) Validate(ba string) error {
	encodedAuth := hash(ba)
	if err := a.DB.CheckAuth(encodedAuth); err != nil {
		return err
	}

	return nil
}

// hash - hashes incoming basic auths encoded in B64 the result will be used check the DB to authorize a user.
func hash(ba string) string {
	h := sha256.New()
	h.Write([]byte(ba))

	return fmt.Sprintf("%x", h.Sum(nil))
}
