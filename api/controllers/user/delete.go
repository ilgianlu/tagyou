package user

import (
	"net/http"

	"github.com/ilgianlu/tagyou/persistence"
	"github.com/rs/zerolog/log"

	"github.com/julienschmidt/httprouter"
)

func (uc UserController) RemoveClient(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user, err := uc.getOne(w, r, p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := persistence.UserRepository.DeleteById(user.ID); err != nil {
		log.Printf("error deleting user: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
