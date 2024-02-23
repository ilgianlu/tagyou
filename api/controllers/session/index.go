package session

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

func (sc SessionController) GetSessions(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sessions := persistence.SessionRepository.GetAll()
	if res, err := json.Marshal(sessions); err != nil {
		slog.Error("error marshaling session rows", "err", err)
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
