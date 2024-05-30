package user

import (
	"log/slog"
	"net/http"

	"github.com/ilgianlu/tagyou/persistence"
)

func (uc UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	user, err := uc.getOne(w, r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := persistence.UserRepository.DeleteById(user.ID); err != nil {
		slog.Error("error deleting user", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
