package subscription

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/julienschmidt/httprouter"
)

func (sc SubscriptionController) GetSubscriptions(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	subscriptions := []model.Subscription{}
	if err := sc.db.Find(&subscriptions).Error; err != nil {
		log.Printf("error getting subscriptions rows: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res, err := json.Marshal(subscriptions); err != nil {
		log.Printf("error marshaling subscriptions rows: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		numBytes, err := w.Write(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		log.Printf("Wrote %d bytes json result\n", numBytes)
	}
}
