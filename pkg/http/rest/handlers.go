package rest

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rnov/Go-REST/pkg/auth"
	mid "github.com/rnov/Go-REST/pkg/http/middleware"
)

type RecipeAPI interface {
	GetRecipeByID(w http.ResponseWriter, r *http.Request)
	GetAllRecipes(w http.ResponseWriter, r *http.Request)
	CreateRecipe(w http.ResponseWriter, r *http.Request)
	UpdateRecipe(w http.ResponseWriter, r *http.Request)
	DeleteRecipe(w http.ResponseWriter, r *http.Request)
}

type RateAPI interface {
	RateRecipe(w http.ResponseWriter, r *http.Request)
}

func NewRouter(rcpHand *RecipeHandler, rateHand *RateHandler, auth *auth.Auth) *mux.Router {
	APIRESTRouter := mux.NewRouter()
	configRecipeEndpoints(APIRESTRouter, rcpHand, auth)
	configRateEndPoints(APIRESTRouter, rateHand)

	return APIRESTRouter
}

func configRecipeEndpoints(r *mux.Router, rcpHand *RecipeHandler, auth *auth.Auth) {
	r.HandleFunc("/recipes/{ID}", rcpHand.GetRecipeByID).Methods("GET")
	r.HandleFunc("/recipes", rcpHand.GetAllRecipes).Methods("GET")
	r.HandleFunc("/recipes/{ID}", mid.Authentication(auth, rcpHand.DeleteRecipe)).Methods("DELETE")
	r.HandleFunc("/recipes", mid.Authentication(auth, rcpHand.CreateRecipe)).Methods("POST")
	r.HandleFunc("/recipes/{ID}", mid.Authentication(auth, rcpHand.UpdateRecipe)).Methods("PUT")
}

func configRateEndPoints(r *mux.Router, rateHand *RateHandler) {
	r.HandleFunc("/recipes/{ID}/rate", rateHand.RateRecipe).Methods("POST")
}
