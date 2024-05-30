package user

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ilgianlu/tagyou/persistence"
)

type IndexUserDTO struct {
	ID        int64
	Username  string
	CreatedAt int64
}

func (uc UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := persistence.UserRepository.GetAll()
	usersDTO := []IndexUserDTO{}
	for _, u := range users {
		usersDTO = append(usersDTO, IndexUserDTO{ID: u.ID, Username: u.Username, CreatedAt: u.CreatedAt})
	}
	if res, err := json.Marshal(usersDTO); err != nil {
		slog.Error("error marshaling auth rows", "err", err)
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
