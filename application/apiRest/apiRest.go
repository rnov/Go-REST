package apiRest

import (
	"Go-REST/application/apiRest/handlers"
	"github.com/julienschmidt/httprouter"
)

func NewRouter(rcpHand *handlers.RecipeHandler, rateHand *handlers.RateHandler) *httprouter.Router {
	ApiRestRouter := &httprouter.Router{}
	configRecipeEndpoints(ApiRestRouter, rcpHand)
	configRateEndPoints(ApiRestRouter, rateHand)

	return ApiRestRouter
}

// note private functions needed to configure route's endpoints, used in NewRouter
func configRecipeEndpoints(router *httprouter.Router, rcpHand *handlers.RecipeHandler) {
	router.GET("/recipes/:id", rcpHand.GetRecipeById)
	router.GET("/recipes", rcpHand.GetAllRecipes)
	router.DELETE("/recipes/:id", rcpHand.DeleteRecipe)
	router.POST("/recipes", rcpHand.CreateRecipe)
	router.PUT("/recipes/:id", rcpHand.UpdateRecipe)
}

func configRateEndPoints(router *httprouter.Router, rateHand *handlers.RateHandler) {
	router.POST("/recipes/:id/rating", rateHand.RateRecipe)
}
