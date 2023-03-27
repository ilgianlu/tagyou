package client

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

type CreateClientDTO struct {
	ClientId             string
	Username             string
	SubscribeAcl         string
	PublishAcl           string
	InputPassword        string
	InputPasswordConfirm string
}

func (a *CreateClientDTO) Validate() bool {
	if a.Username == "" || a.ClientId == "" {
		return false
	}
	if password.ValidPassword([]byte(a.InputPassword)) {
		return false
	}
	return true
}

func (uc ClientController) CreateClient(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	client := CreateClientDTO{}
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		log.Error().Err(err).Msg("error decoding json input")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !client.Validate() {
		log.Error().Msg("data passed is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	pass, _ := password.EncodePassword([]byte(client.InputPassword))
	c := model.Client{
		ClientId:     client.ClientId,
		Username:     client.Username,
		Password:     pass,
		SubscribeAcl: client.SubscribeAcl,
		PublishAcl:   client.PublishAcl,
		CreatedAt:    time.Now().Unix(),
	}

	if err := persistence.ClientRepository.Create(c); err != nil {
		log.Error().Err(err).Msg("error saving new client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
