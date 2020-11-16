package rest

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/rnov/Go-REST/pkg/errors"
	log "github.com/rnov/Go-REST/pkg/logger"
	"github.com/rnov/Go-REST/pkg/rate"
	"github.com/rnov/Go-REST/pkg/service"
)

type RateHandler struct {
	rateSrv service.Rater
	logger  log.Loggers
}

func NewRateHandler(srv service.Rater, l log.Loggers) *RateHandler {
	rateHandler := &RateHandler{
		rateSrv: srv,
		logger:  l,
	}
	return rateHandler
}

func (rh *RateHandler) RateRecipe(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["ID"]
	if len(ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	rating := &rate.Rate{}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(rating); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := rh.rateSrv.Rate(ID, rating); err != nil {
		errors.BuildResponse(w, r.Method, err)
	}

	w.WriteHeader(http.StatusOK)
}
