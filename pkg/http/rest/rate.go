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
	srv    service.Rater
	logger log.Loggers
}

func NewRateHandler(srv service.Rater, l log.Loggers) *RateHandler {
	rateHandler := &RateHandler{
		srv:    srv,
		logger: l,
	}
	return rateHandler
}

func (rh *RateHandler) RateRecipe(w http.ResponseWriter, r *http.Request) {
	ID := mux.Vars(r)["id"]
	if len(ID) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//payload, err := ioutil.ReadAll(r.Body)
	//if err != nil {
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	rating := &rate.Rate{}
	jd := json.NewDecoder(r.Body)
	jd.DisallowUnknownFields()
	if err := jd.Decode(&rating); err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	if err := rh.srv.Rate(ID, rating); err != nil {
		errors.BuildResponse(w, err)
	}

	w.WriteHeader(http.StatusOK)
}
