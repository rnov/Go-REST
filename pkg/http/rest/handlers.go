package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rnov/Go-REST/pkg/auth"
	"github.com/rnov/Go-REST/pkg/errors"
	mid "github.com/rnov/Go-REST/pkg/http/middleware"
	"net/http"
)

type RecipeAPI interface {
	GetRecipeById(w http.ResponseWriter, r *http.Request)
	GetAllRecipes(w http.ResponseWriter, r *http.Request)
	CreateRecipe(w http.ResponseWriter, r *http.Request)
	UpdateRecipe(w http.ResponseWriter, r *http.Request)
	DeleteRecipe(w http.ResponseWriter, r *http.Request)
}

type RateAPI interface {
	RateRecipe(w http.ResponseWriter, r *http.Request)
}

func NewRouter(rcpHand *RecipeHandler, rateHand *RateHandler, auth *auth.Auth) *mux.Router {
	ApiRestRouter := mux.NewRouter()
	configRecipeEndpoints(ApiRestRouter, rcpHand, auth)
	configRateEndPoints(ApiRestRouter, rateHand)

	return ApiRestRouter
}

// note private functions needed to configure route's endpoints, used in NewRouter
func configRecipeEndpoints(r *mux.Router, rcpHand *RecipeHandler, auth *auth.Auth) {

	r.HandleFunc("/recipes/{id}", rcpHand.GetRecipeById).Methods("GET")
	r.HandleFunc("/recipes", rcpHand.GetAllRecipes).Methods("GET")
	r.HandleFunc("/recipes/{id}", mid.Authentication(auth, rcpHand.DeleteRecipe)).Methods("DELETE")
	r.HandleFunc("/recipes", mid.Authentication(auth, rcpHand.CreateRecipe)).Methods("POST")
	r.HandleFunc("/recipes/{id}", mid.Authentication(auth, rcpHand.UpdateRecipe)).Methods("PUT")
}

func configRateEndPoints(r *mux.Router, rateHand *RateHandler) {
	r.HandleFunc("/recipes/{id}/rating", rateHand.RateRecipe).Methods("POST")
}

func BuildErrorResponse(w http.ResponseWriter, err error) {

	switch e := err.(type) {
	case *errors.FailedAuthErr:
		w.WriteHeader(http.StatusUnauthorized)
	case *errors.DBErr:
		//http.Error(http.StatusInternalServerError, ,)
		w.WriteHeader(http.StatusInternalServerError)
	case *errors.NotFoundErr:
		w.WriteHeader(http.StatusNotFound)
	case *errors.InvalidParamsErr:
		body, jsonErr := json.Marshal(e.Parameters)
		if _, parseErr := w.Write(body); jsonErr != nil || parseErr != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
	}
}
