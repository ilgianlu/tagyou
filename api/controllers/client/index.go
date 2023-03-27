package client

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

type IndexClientDTO struct {
	ClientId  string
	Username  string
	CreatedAt int64
}

func (uc ClientController) GetClients(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	clients := persistence.ClientRepository.GetAll()
	clientDTOs := []IndexClientDTO{}
	for _, u := range clients {
		clientDTOs = append(clientDTOs, IndexClientDTO{ClientId: u.ClientId, Username: u.Username, CreatedAt: u.CreatedAt})
	}
	if res, err := json.Marshal(clientDTOs); err != nil {
		log.Printf("error marshaling client rows: %s\n", err)
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
