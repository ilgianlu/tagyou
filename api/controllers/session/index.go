package session

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ilgianlu/tagyou/model"
	"github.com/julienschmidt/httprouter"
)

func (sc SessionController) GetSessions(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sessions := []model.Session{}
	if err := sc.db.Find(&sessions).Error; err != nil {
		log.Printf("error getting session rows: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res, err := json.Marshal(sessions); err != nil {
		log.Printf("error marshaling session rows: %s\n", err)
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
