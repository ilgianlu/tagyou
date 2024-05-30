package subscription

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ilgianlu/tagyou/persistence"
)

func (sc SubscriptionController) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	subscriptions := persistence.SubscriptionRepository.GetAll()
	if res, err := json.Marshal(subscriptions); err != nil {
		slog.Error("error marshaling subscriptions rows", "err", err)
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
		slog.Info("Wrote json result", "num-bytes", numBytes)
	}
}
