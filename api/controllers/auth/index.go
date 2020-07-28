package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ilgianlu/tagyou/model"
	"github.com/julienschmidt/httprouter"
)

func (uc AuthController) GetAuths(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	auths := []model.Auth{}
	if err := uc.db.Find(&auths).Error; err != nil {
		log.Printf("error getting auth rows: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res, err := json.Marshal(auths); err != nil {
		log.Printf("error marshaling auth rows: %s\n", err)
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
