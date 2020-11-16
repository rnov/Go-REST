package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"

	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/logger"
	r "github.com/rnov/Go-REST/pkg/recipe"
)

type RecipeServiceMock struct {
	getByID func(recipeID string) (*r.Recipe, error)
	listAll func() ([]*r.Recipe, error)
	create  func(recipe *r.Recipe) error
	update  func(ID string, recipe *r.Recipe) error
	delete  func(recipeID string) error
}

func (rsm RecipeServiceMock) GetByID(recipeID string) (*r.Recipe, error) {
	if rsm.getByID != nil {
		return rsm.getByID(recipeID)
	}
	panic("Not implemented")
}

func (rsm RecipeServiceMock) ListAll() ([]*r.Recipe, error) {
	if rsm.listAll != nil {
		return rsm.listAll()
	}
	panic("Not implemented")
}

func (rsm RecipeServiceMock) Create(recipe *r.Recipe) error {
	if rsm.create != nil {
		return rsm.create(recipe)
	}
	panic("Not implemented")
}

func (rsm RecipeServiceMock) Update(ID string, recipe *r.Recipe) error {
	if rsm.update != nil {
		return rsm.update(ID, recipe)
	}
	panic("Not implemented")
}

func (rsm RecipeServiceMock) Delete(recipeID string) error {
	if rsm.delete != nil {
		return rsm.delete(recipeID)
	}
	panic("Not implemented")
}

func TestRecipeHandler_GetRecipeByID(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		service         RecipeServiceMock
		status          int
		expectedPayload *r.Recipe
	}{
		{
			name: "Successful request",
			url:  "/recipes/5f10223c",
			service: RecipeServiceMock{
				getByID: func(recipeID string) (*r.Recipe, error) {
					return &r.Recipe{
						ID:         "5f10223c",
						Name:       "qwerty",
						PrepTime:   20,
						Difficulty: 3,
						Vegetarian: false,
					}, nil
				},
			},
			status: 200,
			expectedPayload: &r.Recipe{
				ID:         "5f10223c",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
		},
		{
			name: "error getting recipe ID - out of scope",
			url:  "/recipes/0123456789xyz",
			service: RecipeServiceMock{
				getByID: func(recipeID string) (*r.Recipe, error) {
					return nil, errors.NewInputError("Invalid ID format", nil)
				},
			},
			status: 400,
		},
		{
			name: "error - recipe not found ",
			url:  "/recipes/5f10223c",
			service: RecipeServiceMock{
				getByID: func(recipeID string) (*r.Recipe, error) {
					return nil, errors.NewExistErr(false)
				},
			},
			status: 404,
		},
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
			servicesRouter.HandleFunc("/recipes/{ID}", rh.GetRecipeByID).Methods("GET")

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			servicesRouter.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.status {
				t.Errorf("handler returned wrong status code: expected %v got %v", test.status, rr.Code)
			}

			if test.expectedPayload != nil {
				rcp := &r.Recipe{}
				_ = json.Unmarshal(rr.Body.Bytes(), rcp)
				if !reflect.DeepEqual(test.expectedPayload, rcp) {
					t.Errorf("expected: '%v' instead got: '%v'", test.expectedPayload, rcp)
				}
			}
		})
	}
}

func TestRecipeHandler_GetAllRecipes(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		service         RecipeServiceMock
		status          int
		expectedPayload []*r.Recipe
	}{
		{
			name: "successful request - multiple results",
			url:  "/recipes",
			service: RecipeServiceMock{
				listAll: func() ([]*r.Recipe, error) {
					return []*r.Recipe{
						{

							ID:         "5f10223c",
							Name:       "qwerty",
							PrepTime:   20,
							Difficulty: 3,
							Vegetarian: false,
						},
						{
							ID:         "c32201f5",
							Name:       "ytrewq",
							PrepTime:   25,
							Difficulty: 5,
							Vegetarian: false,
						},
					}, nil
				},
			},
			status: 200,
			expectedPayload: []*r.Recipe{
				{

					ID:         "5f10223c",
					Name:       "qwerty",
					PrepTime:   20,
					Difficulty: 3,
					Vegetarian: false,
				},
				{
					ID:         "c32201f5",
					Name:       "ytrewq",
					PrepTime:   25,
					Difficulty: 5,
					Vegetarian: false,
				},
			},
		},
		{
			name: "successful request - empty result",
			url:  "/recipes",
			service: RecipeServiceMock{
				listAll: func() ([]*r.Recipe, error) {
					return []*r.Recipe{}, nil
				},
			},
			status:          200,
			expectedPayload: nil,
		},
		{
			name: "error - system failure",
			url:  "/recipes",
			service: RecipeServiceMock{
				listAll: func() ([]*r.Recipe, error) {
					return nil, errors.NewDBErr("system failure")
				},
			},
			status: 500,
		},
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

			if test.expectedPayload != nil {
				rcps := make([]*r.Recipe, 0)
				_ = json.Unmarshal(rr.Body.Bytes(), &rcps)
				if !reflect.DeepEqual(test.expectedPayload, rcps) {
					t.Errorf("requestPayload values do not match")
				}
			}
		})
	}
}

func TestRecipeHandler_CreateRecipe(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		requestPayload  *r.Recipe
		service         RecipeServiceMock
		status          int
		expectedPayload *r.Recipe
	}{
		{
			name: "Successful request",
			url:  "/recipes",
			requestPayload: &r.Recipe{
				ID:         "5f10223c",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			service: RecipeServiceMock{
				create: func(recipe *r.Recipe) error {
					return nil
				},
			},
			status: 201,
			expectedPayload: &r.Recipe{
				ID:         "5f10223c",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
		},
		{
			name: "error - invalid recipe params ",
			url:  "/recipes",
			requestPayload: &r.Recipe{
				ID:         "5f10223c",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 100,
				Vegetarian: false,
			},
			service: RecipeServiceMock{
				create: func(recipe *r.Recipe) error {
					return errors.NewInputError("Invalid input parameters", map[string]string{errors.Rate: errors.OutOfRange})
				},
			},
			status: 400,
		},
		{
			name:   "error special case incoming body is not a recipe - error unmarshal",
			url:    "/recipes",
			status: 400,
		},
	}

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

			if test.expectedPayload != nil {
				rcp := &r.Recipe{}
				_ = json.Unmarshal(rr.Body.Bytes(), rcp)
				if !reflect.DeepEqual(test.expectedPayload, rcp) {
					t.Errorf("expected: '%v' instead got: '%v'", test.expectedPayload, rcp)
				}
			}
		})
	}
}

func TestRecipeHandler_UpdateRecipe(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		requestPayload  *r.Recipe
		service         RecipeServiceMock
		status          int
		expectedPayload *r.Recipe
	}{
		{
			name: "Successful request",
			url:  "/recipes/5f10223c",
			requestPayload: &r.Recipe{
				ID:         "5f10223c",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
			service: RecipeServiceMock{
				update: func(ID string, recipe *r.Recipe) error {
					return nil
				},
			},
			status: 200,
			expectedPayload: &r.Recipe{
				ID:         "5f10223c",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
		},
		{
			name: "error - invalid recipe params ",
			url:  "/recipes/5f10223c",
			requestPayload: &r.Recipe{
				ID:         "5f10223c",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 100,
				Vegetarian: false,
			},
			service: RecipeServiceMock{
				update: func(ID string, recipe *r.Recipe) error {
					return errors.NewInputError("Invalid input parameters", map[string]string{errors.Rate: errors.OutOfRange})
				},
			},
			status: 400,
		},
		{
			name: "error updating a recipe that do not exist",
			url:  "/recipes/5f10223c",
			requestPayload: &r.Recipe{
				ID:         "5f10223c",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 100,
				Vegetarian: false,
			},
			service: RecipeServiceMock{
				update: func(ID string, recipe *r.Recipe) error {
					return errors.NewExistErr(false)
				},
			},
			status: 204,
		},
		//{
		//	name:   "error special case incoming body is not a recipe - error unmarshal",
		//	url:    "/recipes/5f10223c",
		//	status: 400,
		//},
	}

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
			req, err := http.NewRequest("PUT", test.url, bytes.NewBuffer(jsonBody))
			if err != nil {
				t.Fatal(err)
			}

			rh := NewRecipeHandler(&test.service, l)

			// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()
			servicesRouter := mux.NewRouter()
			servicesRouter.HandleFunc("/recipes/{ID}", rh.UpdateRecipe).Methods("PUT")

			// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
			// directly and pass in our Request and ResponseRecorder.
			servicesRouter.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if rr.Code != test.status {
				t.Errorf("handler returned wrong status code: expected %v got %v", test.status, rr.Code)
			}

			if test.expectedPayload != nil {
				rcp := &r.Recipe{}
				_ = json.Unmarshal(rr.Body.Bytes(), rcp)
				if !reflect.DeepEqual(test.expectedPayload, rcp) {
					t.Errorf("expected: '%v' instead got: '%v'", test.expectedPayload, rcp)
				}
			}
		})
	}
}

func TestRecipeHandler_DeleteRecipe(t *testing.T) {
	tests := []struct {
		name            string
		url             string
		service         RecipeServiceMock
		status          int
		expectedPayload *r.Recipe
	}{
		{
			name: "Successful request",
			url:  "/recipes/5f10223c",
			service: RecipeServiceMock{
				delete: func(recipeID string) error {
					return nil
				},
			},
			status: 204,
			expectedPayload: &r.Recipe{
				ID:         "5f10223c",
				Name:       "qwerty",
				PrepTime:   20,
				Difficulty: 3,
				Vegetarian: false,
			},
		},
		{
			name: "error getting recipe ID - out of scope",
			url:  "/recipes/0123456789xyz",
			service: RecipeServiceMock{
				delete: func(recipeID string) error {
					return errors.NewInputError("Invalid ID format", nil)
				},
			},
			status: 400,
		},
		{
			name: "error - recipe not found ",
			url:  "/recipes/5f10223c",
			service: RecipeServiceMock{
				delete: func(recipeID string) error {
					return errors.NewExistErr(false)
				},
			},
			status: 404,
		},
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
			servicesRouter.HandleFunc("/recipes/{ID}", rh.DeleteRecipe).Methods("DELETE")

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
