package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/logger"
	"github.com/rnov/Go-REST/pkg/rate"
	"net/http"
	"net/http/httptest"
	"testing"
)

type rateServiceMock struct {
	rate func(id string, r *rate.Rate) error
}

func (rsm *rateServiceMock) Rate(id string, r *rate.Rate) error {
	if rsm.rate != nil {
		return rsm.rate(id, r)
	}
	panic("Not implemented")
}

func TestRateHandler_RateRecipe(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		payload      interface{}
		service      rateServiceMock
		status       int
		checkPayload func(payload string) error
	}{
		//{
		//	name:    "successful rate",
		//	url:     "/recipes/5f10223c-c420-4640-833a-4e4576ca565c/rating",
		//	payload: rate.Rate{Note: 5},
		//	service: rateServiceMock{
		//		rate: func(id string, r *rate.Rate) error {
		//			return nil
		//		},
		//	},
		//	status:       http.StatusOK,
		//	checkPayload: nil,
		//},
		{
			name:         "Invalid Body",
			url:          "/recipes/5f10223c-c420-4640-833a-4e4576ca565c/rating",
			payload:      struct{ name string }{name: "invalid Body"},
			service:      rateServiceMock{},
			status:       http.StatusBadRequest,
			checkPayload: nil,
		},
		{
			name:    "Invalid Data Range",
			url:     "/recipes/5f10223c-c420-4640-833a-4e4576ca565c/rating",
			payload: rate.Rate{Note: 10},
			service: rateServiceMock{
				rate: func(id string, r *rate.Rate) error {
					v := make(map[string]string)
					v[errors.Rate] = errors.OutOfRange
					return errors.NewInvalidParamsErr(v)
				},
			},
			status:       http.StatusBadRequest,
			checkPayload: nil,
		},
	}

	// Create a request to pass to our handler. We don't name have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			l := logger.NewLogger()
			jsonBody, _ := json.Marshal(test.payload)
			req, err := http.NewRequest("POST", test.url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			rh := NewRateHandler(&test.service, l)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/recipes/{id}/rating", rh.RateRecipe).Methods("POST")

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			servicesRouter.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.status {
				t.Errorf("handler returned wrong status code: expected %v got %v", test.status, rr.Code)
			}

			if test.checkPayload != nil {
				if err := test.checkPayload(rr.Body.String()); err != nil {
					t.Errorf("error validation payload: %w", err)
				}
			}

		})
	}
}
