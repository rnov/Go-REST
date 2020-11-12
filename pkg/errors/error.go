package errors

import (
	"encoding/json"
	"net/http"
)

const (
	Exists      = "already Exists"
	NotFound    = "does not Exist"
	AuthFailed  = "auth Failed"
	OutOfRange  = "out of range"
	TooLong     = "too long"
	MissingName = "missing name"
)

const (
	Name       = "name"
	Preptime   = "preptime"
	Difficulty = "difficulty"
)

const (
	Rate   = "rate"
	RateID = "id"
)

const (
	Bearer = "bearer"
)

type ExistErr struct {
	msgToLog string
}

func (myErr *ExistErr) Error() string {
	return myErr.msgToLog
}

func NewExistErr(msg string) *ExistErr {
	return &ExistErr{
		msgToLog: msg,
	}
}

type NotFoundErr struct {
}

func (myErr *NotFoundErr) Error() string {
	return "item not found"
}

func NewNotFoundErr() *NotFoundErr {
	return &NotFoundErr{}
}

type DBErr struct {
	msgToLog string
}

func (myErr *DBErr) Error() string {
	return myErr.msgToLog
}

func NewDBErr(msg string) *DBErr {
	return &DBErr{
		msgToLog: msg,
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

type UserErr struct {
	msg string
}

func (myErr *UserErr) Error() string {
	return myErr.msg
}

func NewUserErr(message string) *UserErr {
	return &UserErr{
		msg: message,
	}
}

func BuildResponse(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *FailedAuthErr:
		w.WriteHeader(http.StatusUnauthorized)
	case *DBErr:
		w.WriteHeader(http.StatusInternalServerError)
	case *NotFoundErr:
		w.WriteHeader(http.StatusNotFound)
	case *UserErr:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(e.Error()))
	case *InvalidParamsErr:
		body, jsonErr := json.Marshal(e.Parameters)
		if jsonErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
	}
}
