package user

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/password"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
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

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := CreateUserDTO{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error().Err(err).Msg("error decoding json input")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !user.Validate() {
		log.Error().Msg("data passed is invalid")
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
		log.Error().Err(err).Msg("error saving new client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
