package rest

import (
	
	"encoding/json"
	"fmt"
	"github.com/goRest/pkg/rates"
	"github.com/goRest/pkg/recipes"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strings"
	
	
)

type RecipeAPI interface {
	GetById(recipeId string) (*recipes.Recipe, error)
	ListAll() ([]*recipes.Recipe, error)
	Create(recipe *recipes.Recipe, auth string) (map[string]string, error)
	Update(recipe *recipes.Recipe, urlId string, auth string) (map[string]string, error)
	Delete(recipeId string, auth string) error
}

type RateAPI interface {
	Rate(id string, rating *rates.Rate) (map[string]string, error) // rate recipe
}

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

type RateHandler struct {
	rateController controller.ApiRateCalls
}

func NewRateHandler(rateInterface controller.ApiRateCalls) *RateHandler {
	rateHandler := &RateHandler{
		rateController: rateInterface,
	}
	return rateHandler
}

func (rateHand *RateHandler) RateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName("id")
	if len(id) == 0 {
		w.WriteHeader(400)
		return
	}

	rating := &rates.Rate{}
	err := json.NewDecoder(r.Body).Decode(&rating)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	invalidParams, err := rateHand.rateController.Rate(id, rating)
	if err != nil {
		if err.Error() == msg.DbError {
			w.WriteHeader(500)
			return
		} else if err.Error() == msg.InvalidParams {
			// Marshal provided interface into JSON structure
			jsonParams, _ := json.Marshal(invalidParams)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			fmt.Fprintf(w, "%s", jsonParams)
			return
		} else if err.Error() == msg.NotFound {
			w.WriteHeader(404)
			return
		}
	}
	w.WriteHeader(200)

}

// interface, could get any controller that implements the interface (redis, mongo, psql ...)
type RecipeHandler struct {
	recipeController controller.ApiRecipeCalls
}

func NewRecipeHandler(rcpInterface controller.ApiRecipeCalls) *RecipeHandler {
	recipeHandler := &RecipeHandler{
		recipeController: rcpInterface,
	}
	return recipeHandler
}

func (rcphand *RecipeHandler) GetRecipeById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName(msg.RecipeId)
	if len(id) == 0 {
		w.WriteHeader(400)
		return
	}

	recipe, err := rcphand.recipeController.GetById(id)

	if err != nil {
		if err.Error() == msg.DbError {
			w.WriteHeader(500)
			return
		} else if err.Error() == msg.NotFound {
			w.WriteHeader(404)
			return
		}
	}

	// Marshal provided interface into JSON structure
	recipeJson, err := json.Marshal(recipe)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", recipeJson)
}

func (rcphand *RecipeHandler) GetAllRecipes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	recipes, err := rcphand.recipeController.ListAll()

	if err != nil {
		if err.Error() == msg.DbError {
			w.WriteHeader(500)
			return
		}
	}

	var recipesJson []byte
	if len(recipes) == 0 {
		recipesJson, err = json.Marshal([]string{})
	} else {
		recipesJson, err = json.Marshal(recipes)
	}

	if err != nil {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", recipesJson)
}

func (rcphand *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	recipe := recipes.Recipe{}
	err := json.NewDecoder(r.Body).Decode(&recipe)

	if err != nil {
		w.WriteHeader(400)
		return
	}

	auth := r.Header[msg.Authentication]
	if len(auth) != 1 {
		w.WriteHeader(401)
		return
	}
	auth = strings.Fields(auth[0])
	if validateAuthStructure(auth) {
		w.WriteHeader(401)
		return
	}

	invalidParams, err := rcphand.recipeController.Create(&recipe, auth[1])

	if err != nil {

		if err.Error() == msg.AuthFailed {
			w.WriteHeader(401)
			return
		} else if err.Error() == msg.DbError {
			w.WriteHeader(500)
			return
		} else if err.Error() == msg.Exists {
			w.WriteHeader(409)
			return
		} else if err.Error() == msg.InvalidParams {
			jsonParams, err := json.Marshal(invalidParams)
			if err != nil {
				w.WriteHeader(500)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			fmt.Fprintf(w, "%s", jsonParams)
			return
		}


	}

	w.WriteHeader(201)
}

func (rcphand *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	recipe := recipes.Recipe{}
	err := json.NewDecoder(r.Body).Decode(&recipe)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	auth := r.Header[msg.Authentication]
	if len(auth) != 1 {
		w.WriteHeader(400)
		return
	}
	auth = strings.Fields(auth[0])
	if validateAuthStructure(auth) {
		w.WriteHeader(401)
		return
	}

	id := p.ByName(msg.RecipeId)
	if len(id) == 0 {
		w.WriteHeader(400)
		return
	}

	invalidParams, err := rcphand.recipeController.Update(&recipe, id, auth[1])
	if err != nil {
		if err.Error() == msg.AuthFailed {
			w.WriteHeader(401)
			return
		} else if err.Error() == msg.DbError {
			w.WriteHeader(500)
			return
		} else if err.Error() == msg.NotFound {
			w.WriteHeader(404)
			return
		} else if err.Error() == msg.InvalidParams {
			jsonParams, _ := json.Marshal(invalidParams)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			fmt.Fprintf(w, "%s", jsonParams)
			return
		}
	}

	w.WriteHeader(200)
}

func (rcphand *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName(msg.RecipeId)
	if len(id) == 0 {
		w.WriteHeader(400)
		return
	}

	auth := r.Header[msg.Authentication]
	if len(auth) != 1 {
		w.WriteHeader(401)
		return
	}
	auth = strings.Fields(auth[0])
	if validateAuthStructure(auth) {
		w.WriteHeader(401)
		return
	}

	err := rcphand.recipeController.Delete(id, auth[1])
	if err != nil {
		if err.Error() == msg.DbError {
			w.WriteHeader(500)
			return
		} else if err.Error() == msg.AuthFailed {
			w.WriteHeader(401)
			return
		} else if err.Error() == msg.NotFound {
			w.WriteHeader(404)
			return
		}
	}

	w.WriteHeader(204)
}

func validateAuthStructure(auth []string) bool {
	return len(auth) != 2 || auth[0] != msg.Bearer
}
