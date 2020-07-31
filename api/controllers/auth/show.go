package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (uc AuthController) GetAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	auth, err := uc.getOne(w, r, p)
	if err != nil {
		return
	}

	if res, err := json.Marshal(auth); err != nil {
		log.Printf("error marshaling auth row: %s\n", err)
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
