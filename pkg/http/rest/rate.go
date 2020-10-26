package rest

import (
	"encoding/json"
	log "github.com/rnov/Go-REST/pkg/logger"
	"github.com/rnov/Go-REST/pkg/rate"
	"net/http"
)

type RateHandler struct {
	srv    rate.Rater
	logger log.Loggers
}

func NewRateHandler(srv rate.Rater, l log.Loggers) *RateHandler {
	rateHandler := &RateHandler{
		srv:    srv,
		logger: l,
	}
	return rateHandler
}

func (rh *RateHandler) RateRecipe(w http.ResponseWriter, r *http.Request) {

	ID := r.Header.Get("ID")
	if len(ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rating := &rate.Rate{}
	if err := json.NewDecoder(r.Body).Decode(&rating); err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	if err := rh.srv.Rate(ID, rating); err != nil {
		BuildErrorResponse(w, err)
	}

	w.WriteHeader(http.StatusOK)
}
