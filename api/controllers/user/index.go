package user

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	users := persistence.UserRepository.GetAll()
	if res, err := json.Marshal(users); err != nil {
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
