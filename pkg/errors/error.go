package errors

import (
	"encoding/json"
	"net/http"
)

const (
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
	RateID = "ID"
)

// DBErr is a defined error type whose purpose is to be used whenever a DB related error has occurred and needs to be
//logged.
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

// FailedAuthErr is a defined error type whose purpose is to acknowledge a failed authorization attempt.
type FailedAuthErr struct {
}

func (myErr *FailedAuthErr) Error() string {
	return "Failed validation attempt "
}

func NewFailedAuthErr() *FailedAuthErr {
	return &FailedAuthErr{}
}

// InputErr is a defined error type whose purpose is to acknowledge any error due user's input and carry relevant
//information regarding the error.
type InputErr struct {
	Msg        string            `json:"error,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

func (ie *InputErr) Error() string {
	return ie.Msg
}

func NewInputError(msg string, params map[string]string) *InputErr {
	return &InputErr{
		Msg:        msg,
		Parameters: params,
	}
}

// ExistErr is a defined error type whose purpose is to acknowledge that an error has occurred because the data already
//exists or is produce because is missing.
type ExistErr struct {
	Exist bool
}

func (ee *ExistErr) Error() string {
	var retMsg string
	if ee.Exist {
		retMsg = "item already exists"
	} else {
		retMsg = "item does not Exist"
	}
	return retMsg
}

func NewExistErr(exist bool) *ExistErr {
	return &ExistErr{
		Exist: exist,
	}
}

// BuildResponse - is a method that is being used by handler methods whenever an error occurs and a response based on the error type
// needs to be built. It also acknowledge whether an error needs to be logged, due the logging policy design.
// A compromise decision that tights the relation error-log but for the current size is a small one.
func BuildResponse(w http.ResponseWriter, method string, err error) (toLog bool) {
	switch e := err.(type) {
	case *FailedAuthErr:
		w.WriteHeader(http.StatusUnauthorized)
	case *DBErr:
		toLog = true
		w.WriteHeader(http.StatusInternalServerError)
	case *ExistErr:
		if method == "GET" && !e.Exist {
			w.WriteHeader(http.StatusNotFound)
		} else if method == "POST" && e.Exist {
			w.WriteHeader(http.StatusForbidden)
		} else if method == "PUT" && !e.Exist {
			w.WriteHeader(http.StatusNoContent)
		} else if method == "DELETE" && !e.Exist {
			w.WriteHeader(http.StatusNotFound)
		} else if method == "POST" && !e.Exist {
			w.WriteHeader(http.StatusNotFound)
		}
	case *InputErr:
		body, jsonErr := json.Marshal(e)
		if jsonErr != nil {
			toLog = true
			w.WriteHeader(http.StatusInternalServerError)
			return toLog
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
	default:
		toLog = true
		w.WriteHeader(http.StatusInternalServerError)
	}

	return toLog
}
