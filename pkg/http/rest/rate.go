package rest

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/rnov/Go-REST/pkg/rate"
	"net/http"
)

type RateHandler struct {
	srv rate.Rater
}

func NewRateHandler(srv rate.Rater) *RateHandler {
	rateHandler := &RateHandler{
		srv: srv,
	}
	return rateHandler
}

func (rh *RateHandler) RateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName("id")
	if len(id) == 0 {
		w.WriteHeader(400)
		return
	}

	rating := &rate.Rate{}
	err := json.NewDecoder(r.Body).Decode(&rating)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	if err = rh.srv.Rate(id, rating); err != nil{
		buildErrorResponse(w, err)
	}

	w.WriteHeader(200)
}
