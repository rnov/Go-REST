package auth

import (
	"golang.org/x/crypto/bcrypt"

	"github.com/rnov/Go-REST/pkg/db"
	"github.com/rnov/Go-REST/pkg/logger"
)

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

func (a Auth) Validate(ba string) error {
	encodedAuth, err := encode(ba)
	if err != nil {
		return err
	}
	if err := a.DB.CheckAuth(encodedAuth); err != nil {
		return err
	}
	return nil
}

func encode(ba string) (string, error) {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(ba), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	} // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash), nil
}
