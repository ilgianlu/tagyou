package client

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

type IndexClientDTO struct {
	ID        int64
	ClientId  string
	Username  string
	CreatedAt int64
}

func (uc ClientController) GetClients(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	clients := persistence.ClientRepository.GetAll()
	clientDTOs := []IndexClientDTO{}
	for _, u := range clients {
		clientDTOs = append(clientDTOs, IndexClientDTO{ID: u.ID, ClientId: u.ClientId, Username: u.Username, CreatedAt: u.CreatedAt})
	}
	if res, err := json.Marshal(clientDTOs); err != nil {
		slog.Error("error marshaling client rows", "err", err)
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
