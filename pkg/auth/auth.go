package auth

import (
	"log"

	"golang.org/x/crypto/bcrypt"

	"github.com/rnov/Go-REST/pkg/db"
)

type Auth struct {
	DB db.Auth
}

func NewAuth(db db.Auth) *Auth {
	return &Auth{
		DB: db,
	}
}

func (a Auth) Validate(ba string) error {
	encodedAuth := encode(ba)
	if err := a.DB.CheckAuth(encodedAuth); err != nil {
		return err
	}
	return nil
}

func encode(ba string) string {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword([]byte(ba), bcrypt.DefaultCost)
	if err != nil {
		log.Println(err)
	} // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}
