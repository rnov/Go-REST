package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rnov/Go-REST/pkg/errors"
	"github.com/rnov/Go-REST/pkg/logger"
	"github.com/rnov/Go-REST/pkg/recipe"
	"github.com/rnov/Go-REST/pkg/service"
)

const (
	recipeID = "ID"
)

// interface, could get any controller that implements the interface (redis, mongo, psql ...)
type RecipeHandler struct {
	rcpSrv service.RecipeMng
	log    logger.Loggers
	// add a log ? be able to log at handler level ?? move from service and log in here, good idea ?
}

func NewRecipeHandler(rcpSrv service.RecipeMng, l logger.Loggers) *RecipeHandler {
	recipeHandler := &RecipeHandler{
		rcpSrv: rcpSrv,
		log:    l,
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
		errors.BuildResponse(w, r.Method, err)
		return
	}

	// Marshal provided interface into JSON structure
	recipeJSON, err := json.Marshal(rcp)
	if err != nil {
		rh.log.Errorf("system error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Write content-type, status code, requestPayload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", recipeJSON)
}

func (rh *RecipeHandler) GetAllRecipes(w http.ResponseWriter, r *http.Request) {
	rcps, err := rh.rcpSrv.ListAll()
	if err != nil {
		rh.log.Errorf("system error: %s", err.Error())
		errors.BuildResponse(w, r.Method, err)
	}
	var recipesJSON []byte
	w.Header().Set("Content-Type", "application/json")
	recipesJSON, jsonErr := json.Marshal(rcps)
	if jsonErr != nil {
		rh.log.Errorf("system error: %s", jsonErr.Error())
		return
	}
	if _, parseErr := w.Write(recipesJSON); parseErr != nil {
		rh.log.Errorf("system error: %s", parseErr.Error())
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
		errors.BuildResponse(w, r.Method, err)
		return
	}

	body, jsonErr := json.Marshal(rcp)
	if jsonErr != nil {
		rh.log.Errorf("system error: %s", jsonErr.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
		errors.BuildResponse(w, r.Method, err)
		return
	}

	body, err := json.Marshal(rcp)
	if err != nil {
		rh.log.Errorf("system error: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
		errors.BuildResponse(w, r.Method, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
