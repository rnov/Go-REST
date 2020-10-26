package rest

import (
	"encoding/json"
	"fmt"
	log "github.com/rnov/Go-REST/pkg/logger"
	"github.com/rnov/Go-REST/pkg/recipe"
	"net/http"
)

// interface, could get any controller that implements the interface (redis, mongo, psql ...)
type RecipeHandler struct {
	rcpSrv recipe.RcpSrv
	logger log.Loggers
	// add a logger ? be able to log at handler level ?? move from service and log in here, good idea ?
}

func NewRecipeHandler(rcpSrv recipe.RcpSrv, l log.Loggers) *RecipeHandler {
	recipeHandler := &RecipeHandler{
		rcpSrv: rcpSrv,
		logger: l,
	}
	return recipeHandler
}

func (rh *RecipeHandler) GetRecipeById(w http.ResponseWriter, r *http.Request) {

	id := r.Header.Get("msg.RecipeId")
	if len(id) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rcp, err := rh.rcpSrv.GetById(id)
	if err != nil {
		BuildErrorResponse(w, err)
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

func (rh *RecipeHandler) GetAllRecipes(w http.ResponseWriter, r *http.Request) {

	rcps, err := rh.rcpSrv.ListAll()
	if err != nil {
		BuildErrorResponse(w, err)
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

func (rh *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {

	recipe := recipe.Recipe{}
	err := json.NewDecoder(r.Body).Decode(&recipe)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//auth := r.Header["msg.Authentication"]
	//if len(auth) != 1 {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	rcp, err := rh.rcpSrv.Create(&recipe)
	if err != nil {
		BuildErrorResponse(w, err)
	}

	body, jsonErr := json.Marshal(rcp)
	if _, parseErr := w.Write(body); jsonErr != nil || parseErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (rh *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
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
	ID := r.Header.Get("msg.RecipeId")
	if len(ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = rh.rcpSrv.Update(&recipe, ID, auth[1])
	if err != nil {
		BuildErrorResponse(w, err)
	}

	w.WriteHeader(http.StatusOK)
}

func (rh *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {

	ID := r.Header.Get("RecipeId")
	if len(ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := rh.rcpSrv.Delete(ID); err != nil {
		BuildErrorResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
