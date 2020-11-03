package rest

import (
	"github.com/gorilla/mux"
	"github.com/rnov/Go-REST/pkg/auth"
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


