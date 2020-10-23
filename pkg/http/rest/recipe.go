package rest

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/rnov/Go-REST/pkg/recipe"
	"net/http"
	"strings"
)

// interface, could get any controller that implements the interface (redis, mongo, psql ...)
type RecipeHandler struct {
	rcpSrv recipe.RcpSrv
	// add a logger ? be able to log at handler level ?? move from service and log in here, good idea ?
}

func NewRecipeHandler(rcpSrv recipe.RcpSrv) *RecipeHandler {
	recipeHandler := &RecipeHandler{
		rcpSrv: rcpSrv,
	}
	return recipeHandler
}

func (rh *RecipeHandler) GetRecipeById(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName("msg.RecipeId")
	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rcp, err := rh.rcpSrv.GetById(id)
	if err != nil {
		buildErrorResponse(w, err)
	}

	// Marshal provided interface into JSON structure
	recipeJson, err := json.Marshal(rcp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write content-type, status code, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", recipeJson)
}

func (rh *RecipeHandler) GetAllRecipes(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	rcps, err := rh.rcpSrv.ListAll()
	if err != nil {
		buildErrorResponse(w, err)
	}

	var recipesJson []byte
	w.Header().Set("Content-Type", "application/json")
	recipesJson, jsonErr := json.Marshal(rcps)
	w.WriteHeader(http.StatusOK)
	if _, parseErr := w.Write(recipesJson); jsonErr != nil || parseErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rh *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	recipe := recipe.Recipe{}
	err := json.NewDecoder(r.Body).Decode(&recipe)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth := r.Header["msg.Authentication"]
	if len(auth) != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	auth = strings.Fields(auth[0])
	if validateAuthStructure(auth) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	rcp, err := rh.rcpSrv.Create(&recipe, auth[1])
	if err != nil {
		buildErrorResponse(w, err)
	}

	body, jsonErr := json.Marshal(rcp)
	if _, parseErr := w.Write(body); jsonErr != nil || parseErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rh *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	recipe := recipe.Recipe{}
	err := json.NewDecoder(r.Body).Decode(&recipe)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	auth := r.Header["msg.Authentication"]
	if len(auth) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	auth = strings.Fields(auth[0])
	if validateAuthStructure(auth) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	id := p.ByName("msg.RecipeId")
	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rh.rcpSrv.Update(&recipe, id, auth[1])
	if err != nil {
		buildErrorResponse(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (rh *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName("RecipeId")
	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	auth := r.Header["Authentication"]
	if len(auth) != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	auth = strings.Fields(auth[0])
	if validateAuthStructure(auth) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := rh.rcpSrv.Delete(id, auth[1]); err != nil {
		buildErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
