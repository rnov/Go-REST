package rest

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rnov/Go-REST/pkg/errors"
	log "github.com/rnov/Go-REST/pkg/logger"
	"github.com/rnov/Go-REST/pkg/recipe"
	"github.com/rnov/Go-REST/pkg/service"
	"net/http"
)

const (
	recipeID = "ID"
)

// interface, could get any controller that implements the interface (redis, mongo, psql ...)
type RecipeHandler struct {
	rcpSrv service.RcpSrv
	logger log.Loggers
	// add a logger ? be able to log at handler level ?? move from service and log in here, good idea ?
}

func NewRecipeHandler(rcpSrv service.RcpSrv, l log.Loggers) *RecipeHandler {
	recipeHandler := &RecipeHandler{
		rcpSrv: rcpSrv,
		logger: l,
	}
	return recipeHandler
}

func (rh *RecipeHandler) GetRecipeByID(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	ID := params[recipeID]
	if len(ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rcp, err := rh.rcpSrv.GetByID(ID)
	if err != nil {
		errors.BuildResponse(w, err)
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
		errors.BuildResponse(w, err)
	}
	var recipesJson []byte
	w.Header().Set("Content-Type", "application/json")
	recipesJson, jsonErr := json.Marshal(rcps)
	if _, parseErr := w.Write(recipesJson); jsonErr != nil || parseErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (rh *RecipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {

	rcp := &recipe.Recipe{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(rcp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := rh.rcpSrv.Create(rcp); err != nil {
		errors.BuildResponse(w, err)
		return
	}

	body, jsonErr := json.Marshal(rcp)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (rh *RecipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	rcp := &recipe.Recipe{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(rcp); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	params := mux.Vars(r)
	ID := params[recipeID]
	if len(ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := rh.rcpSrv.Update(ID, rcp); err != nil {
		errors.BuildResponse(w, err)
	}

	body, jsonErr := json.Marshal(rcp)
	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (rh *RecipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ID := params[recipeID]
	if len(ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := rh.rcpSrv.Delete(ID); err != nil {
		errors.BuildResponse(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
