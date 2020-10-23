package rest

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/rnov/Go-REST/pkg/errors"
	"net/http"
)

type RecipeAPI interface {
	GetRecipeById(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	GetAllRecipes(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	CreateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	UpdateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params)
	DeleteRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

type RateAPI interface {
	RateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params)
}

func NewRouter(rcpHand *RecipeHandler, rateHand *RateHandler) *httprouter.Router {
	ApiRestRouter := &httprouter.Router{}
	configRecipeEndpoints(ApiRestRouter, rcpHand)
	configRateEndPoints(ApiRestRouter, rateHand)

	return ApiRestRouter
}

// note private functions needed to configure route's endpoints, used in NewRouter
func configRecipeEndpoints(router *httprouter.Router, rcpHand *RecipeHandler) {
	router.GET("/recipes/:id", rcpHand.GetRecipeById)
	router.GET("/recipes", rcpHand.GetAllRecipes)
	router.DELETE("/recipes/:id", rcpHand.DeleteRecipe)
	router.POST("/recipes", rcpHand.CreateRecipe)
	router.PUT("/recipes/:id", rcpHand.UpdateRecipe)
}

func configRateEndPoints(router *httprouter.Router, rateHand *RateHandler) {
	router.POST("/recipes/:id/rating", rateHand.RateRecipe)
}

func validateAuthStructure(auth []string) bool {
	return len(auth) != 2 || auth[0] != "msg.Bearer"
}

func buildErrorResponse(w http.ResponseWriter, err error) {

	switch e := err.(type) {
	case *errors.AuthFailedErr:
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
