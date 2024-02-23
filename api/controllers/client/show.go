package client

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type ShowClientDTO struct {
	ID        uint
	ClientId  string
	Username  string
	CreatedAt int64
}

func (uc ClientController) GetClient(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	client, err := uc.getOne(w, r, p)
	if err != nil {
		return
	}
	clientDTO := ShowClientDTO{
		ID:        client.ID,
		ClientId:  client.ClientId,
		Username:  client.Username,
		CreatedAt: client.CreatedAt,
	}

	if res, err := json.Marshal(clientDTO); err != nil {
		slog.Error("error marshaling client row", "err", err)
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
