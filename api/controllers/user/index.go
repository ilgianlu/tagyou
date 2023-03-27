package user

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

type IndexUserDTO struct {
	ID        uint
	Username  string
	CreatedAt int64
}

func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	users := persistence.UserRepository.GetAll()
	usersDTO := []IndexUserDTO{}
	for _, u := range users {
		usersDTO = append(usersDTO, IndexUserDTO{ID: u.ID, Username: u.Username, CreatedAt: u.CreatedAt})
	}
	if res, err := json.Marshal(usersDTO); err != nil {
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
