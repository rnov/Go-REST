package rest

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rnov/Go-REST/pkg/errors"
	log "github.com/rnov/Go-REST/pkg/logger"
	"github.com/rnov/Go-REST/pkg/rate"
	"github.com/rnov/Go-REST/pkg/service"
	"net/http"
)

type RateHandler struct {
	rateSrv service.RateSrv
	logger  log.Loggers
}

func NewRateHandler(srv service.RateSrv, l log.Loggers) *RateHandler {
	rateHandler := &RateHandler{
		rateSrv: srv,
		logger:  l,
	}
	return rateHandler
}

func (rh *RateHandler) RateRecipe(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]
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
