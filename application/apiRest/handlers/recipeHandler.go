package handlers

import (
	msg "Go-REST/application/common"
	"Go-REST/application/controller"
	"Go-REST/application/model"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

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

	recipe := model.Recipe{}
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
	recipe := model.Recipe{}
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
