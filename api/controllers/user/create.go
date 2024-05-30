package user

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/password"
	"github.com/ilgianlu/tagyou/persistence"
)

type CreateUserDTO struct {
	Username             string
	InputPassword        string
	InputPasswordConfirm string
}

func (a *CreateUserDTO) Validate() bool {
	if a.Username == "" {
		return false
	}
	if password.ValidPassword([]byte(a.InputPassword)) {
		return false
	}
	return true
}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	user := CreateUserDTO{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("error decoding json input", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !user.Validate() {
		slog.Warn("data passed is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pass, _ := password.EncodePassword([]byte(user.InputPassword))
	u := model.User{
		Username:  user.Username,
		Password:  pass,
		CreatedAt: time.Now().Unix(),
	}
	if err := persistence.UserRepository.Create(u); err != nil {
		slog.Error("error saving new client", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
