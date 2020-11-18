package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"

	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/logger"
	"github.com/rnov/Go-REST/pkg/rate"
)

type rateServiceMock struct {
	rate func(ID string, r *rate.Rate) error
}

func (rsm *rateServiceMock) Rate(ID string, r *rate.Rate) error {
	if rsm.rate != nil {
		return rsm.rate(ID, r)
	}
	panic("Not implemented")
}

func TestRateHandler_RateRecipe(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		requestPayload *rate.Rate
		service        rateServiceMock
		status         int
	}{
		{
			name:           "successful rate",
			url:            "/recipes/5f10223c/rate",
			requestPayload: &rate.Rate{Note: 5},
			service: rateServiceMock{
				rate: func(ID string, r *rate.Rate) error {
					return nil
				},
			},
			status: http.StatusOK,
		},
		{
			name:           "error - invalid Data Range",
			url:            "/recipes/5f10223c/rate",
			requestPayload: &rate.Rate{Note: 10},
			service: rateServiceMock{
				rate: func(ID string, r *rate.Rate) error {
					v := make(map[string]string)
					v[errors.Rate] = errors.OutOfRange
					return errors.NewInputError("invalid input parameters", v)
				},
			},
			status: http.StatusBadRequest,
		},
	}

	// Create a request to pass to our handler. We don't name have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			l := logger.NewLogger()
			var jsonBody []byte
			if test.requestPayload == nil {
				invalidBody := struct {
					Age int `json:"age"`
				}{
					Age: 100,
				}
				jsonBody, _ = json.Marshal(&invalidBody)
			} else {
				jsonBody, _ = json.Marshal(test.requestPayload)
			}
			jsonBody, _ = json.Marshal(test.requestPayload)
			req, err := http.NewRequest("POST", test.url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			rh := NewRateHandler(&test.service, l)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/recipes/{ID}/rate", rh.RateRecipe).Methods("POST")

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			servicesRouter.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.status {
				t.Errorf("handler returned wrong status code: expected %v got %v", test.status, rr.Code)
			}
		})
	}
}
