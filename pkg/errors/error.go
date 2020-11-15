package errors

import (
	"encoding/json"
	"fmt"
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
}

func (myErr *FailedAuthErr) Error() string {
	return fmt.Sprintf("Failed validation atempt ")
}

func NewFailedAuthErr() *FailedAuthErr {
	return &FailedAuthErr{
	}
}

// Error for user Invalid input parameterss - such errors are not being logged their purpose is to inform user of the missing/wrong data
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

type ExistErr struct {
	exist bool
}

func (ee *ExistErr) Error() string {
	var retMsg string
	if ee.exist {
		retMsg = "item already exists"
	} else {
		retMsg = "item does not exist"
	}
	return retMsg
}

func NewExistErr(exist bool) *ExistErr {
	return &ExistErr{
		exist: exist,
	}
}

func BuildResponse(w http.ResponseWriter, method string, err error) {

	switch e := err.(type) {
	case *FailedAuthErr:
		w.WriteHeader(http.StatusUnauthorized)
	case *DBErr:
		w.WriteHeader(http.StatusInternalServerError)
	case *ExistErr:
		if method == "GET" && !e.exist {
			w.WriteHeader(http.StatusNotFound)
		} else if method == "POST" && e.exist {
			w.WriteHeader(http.StatusForbidden)
		} else if method == "PUT" && !e.exist {
			w.WriteHeader(http.StatusNoContent)
		} else if method == "DELETE" && !e.exist {
			w.WriteHeader(http.StatusNotFound)
		}
	case *InputErr:
		body, jsonErr := json.Marshal(e)
		if jsonErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(body)
	}
}
