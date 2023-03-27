package client

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/julienschmidt/httprouter"
)

func (uc ClientController) GetClient(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	client, err := uc.getOne(w, r, p)
	if err != nil {
		return
	}

	if res, err := json.Marshal(client); err != nil {
		log.Printf("error marshaling client row: %s\n", err)
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
