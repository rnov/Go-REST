package handlers

import (
	msg "Go-REST/application/common"
	"Go-REST/application/controller"
	"Go-REST/application/model"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RateHandler struct {
	rateController controller.ApiRateCalls
}

func NewRateHandler(rateInterface controller.ApiRateCalls) *RateHandler {
	rateHandler := &RateHandler{
		rateController: rateInterface,
	}
	return rateHandler
}

func (rateHand *RateHandler) RateRecipe(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	id := p.ByName("id")
	if len(id) == 0 {
		w.WriteHeader(400)
		return
	}

	rating := &model.Rate{}
	err := json.NewDecoder(r.Body).Decode(&rating)

	if err != nil {
		w.WriteHeader(404)
		return
	}

	invalidParams, err := rateHand.rateController.Rate(id, rating)
	if err != nil {
		if err.Error() == msg.DbError {
			w.WriteHeader(500)
			return
		} else if err.Error() == msg.InvalidParams {
			// Marshal provided interface into JSON structure
			jsonParams, _ := json.Marshal(invalidParams)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			fmt.Fprintf(w, "%s", jsonParams)
			return
		} else if err.Error() == msg.NotFound {
			w.WriteHeader(404)
			return
		}
	}
	w.WriteHeader(200)

}
