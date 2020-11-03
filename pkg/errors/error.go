package errors

import (
	"encoding/json"
	"net/http"
)

const (
	Exists        = "already Exists"
	NotFound      = "does not Exist"
	DBError       = "connection Failed"
	AuthFailed    = "auth Failed"
	InvalidParams = "no valid parameters"

	OutOfRange = "out of range"
	TooLong    = "too long"
)

const (
	Id         = "id"
	Name       = "name"
	Preptime   = "preptime"
	Difficulty = "difficulty"
	Vegetarian = "vegetarian"
)

const (
	Rate   = "rate"
	RateId = "id"
)

const (
	Authentication = "Authentication"
	Bearer         = "bearer"
	RecipeId       = "id"
)

type ExistErr struct {
	msg string
}

func (myErr *ExistErr) Error() string {
	return Exists
}

func NewExistErr(message string) *ExistErr {
	return &ExistErr{
		msg: message,
	}
}

type NotFoundErr struct {
	msg string
}

func (myErr *NotFoundErr) Error() string {
	return NotFound
}

func NewNotFoundErr(message string) *NotFoundErr {
	return &NotFoundErr{
		msg: message,
	}
}

type DBErr struct {
	msg      string
	msgToLog string
}

func (myErr *DBErr) Error() string {
	return DBError
}

func NewDBErr(message string) *DBErr {
	//message := fmt.Sprintf("error in DB")
	return &DBErr{
		msg: message,
	}
}

type FailedAuthErr struct {
	msg string
}

func (myErr *FailedAuthErr) Error() string {
	return AuthFailed
}

func NewFailedAuthErr(message string) *FailedAuthErr {
	return &FailedAuthErr{
		msg: message,
	}
}

type InvalidParamsErr struct {
	msg        string
	Parameters map[string]string
}

func (myErr *InvalidParamsErr) Error() string {
	return "invalid parameters"
}

func NewInvalidParamsErr(params map[string]string) *InvalidParamsErr {
	return &InvalidParamsErr{
		Parameters: params,
	}
}

func BuildErrorResponse(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *FailedAuthErr:
		w.WriteHeader(http.StatusUnauthorized)
	case *DBErr:
		//http.Error(http.StatusInternalServerError, ,)
		w.WriteHeader(http.StatusInternalServerError)
	case *NotFoundErr:
		w.WriteHeader(http.StatusNotFound)
	case *InvalidParamsErr:
		body, jsonErr := json.Marshal(e.Parameters)
		if _, parseErr := w.Write(body); jsonErr != nil || parseErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	}
}
