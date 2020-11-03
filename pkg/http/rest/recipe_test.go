package rest

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rnov/Go-REST/pkg/logger"
	r "github.com/rnov/Go-REST/pkg/recipe"
	"net/http"
	"net/http/httptest"
	"testing"
)

type RecipeServiceMock struct {
	getById func(recipeId string) (*r.Recipe, error)
	listAll func() ([]*r.Recipe, error)
	create  func(recipe *r.Recipe) (*r.Recipe, error)
	update  func(recipe *r.Recipe) error
	delete  func(recipeId string) error
}

func (rsm RecipeServiceMock) GetById(recipeId string) (*r.Recipe, error) {
	if rsm.getById != nil {
		return rsm.getById(recipeId)
	}
	panic("Not implemented")
}

func (rsm RecipeServiceMock) ListAll() ([]*r.Recipe, error) {
	if rsm.listAll != nil {
		return rsm.listAll()
	}
	panic("Not implemented")
}

func (rsm RecipeServiceMock) Create(recipe *r.Recipe) (*r.Recipe, error) {
	if rsm.create != nil {
		return rsm.create(recipe)
	}
	panic("Not implemented")
}

func (rsm RecipeServiceMock) Update(recipe *r.Recipe) error {
	if rsm.create != nil {
		return rsm.update(recipe)
	}
	panic("Not implemented")
}

func (rsm RecipeServiceMock) Delete(recipeId string) error {
	if rsm.delete != nil {
		return rsm.delete(recipeId)
	}
	panic("Not implemented")
}

func TestRecipeHandler_GetRecipeById(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		payload      interface{}
		service      RecipeServiceMock
		status       int
		checkPayload func(payload string) error
	}{
		// todo
		{},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			l := logger.NewLogger()
			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rh := NewRecipeHandler(&test.service, l)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/recipes/{id}", rh.GetRecipeById).Methods("GET")

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

func TestRecipeHandler_GetAllRecipes(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		payload      interface{}
		service      RecipeServiceMock
		status       int
		checkPayload func(payload string) error
	}{
		// todo
		{},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			l := logger.NewLogger()
			req, err := http.NewRequest("GET", test.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rh := NewRecipeHandler(&test.service, l)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/recipes", rh.GetAllRecipes).Methods("GET")

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

func TestRecipeHandler_CreateRecipe(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		payload      interface{}
		service      RecipeServiceMock
		status       int
		checkPayload func(payload string) error
	}{
		// todo
		{},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			l := logger.NewLogger()
			jsonBody, _ := json.Marshal(test.payload)
			req, err := http.NewRequest("POST", test.url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			rh := NewRecipeHandler(&test.service, l)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/recipes", rh.CreateRecipe).Methods("POST")

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

func TestRecipeHandler_UpdateRecipe(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		payload      interface{}
		service      RecipeServiceMock
		status       int
		checkPayload func(payload string) error
	}{
		// todo
		{},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			l := logger.NewLogger()
			jsonBody, _ := json.Marshal(test.payload)
			req, err := http.NewRequest("PUT", test.url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			rh := NewRecipeHandler(&test.service, l)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/recipes/{id}", rh.UpdateRecipe).Methods("PUT")

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

func TestRecipeHandler_DeleteRecipe(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		payload      interface{}
		service      RecipeServiceMock
		status       int
		checkPayload func(payload string) error
	}{
		// todo
		{},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			l := logger.NewLogger()
			req, err := http.NewRequest("DELETE", test.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rh := NewRecipeHandler(&test.service, l)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/recipes/{id}", rh.DeleteRecipe).Methods("DELETE")

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
