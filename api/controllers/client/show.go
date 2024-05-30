package client

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ShowClientDTO struct {
	ID        int64
	ClientId  string
	Username  string
	CreatedAt int64
}

func (uc ClientController) GetClient(w http.ResponseWriter, r *http.Request) {
	client, err := uc.getOne(w, r)
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
