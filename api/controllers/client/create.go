package client

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/password"
	"github.com/ilgianlu/tagyou/persistence"
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

func (uc ClientController) CreateClient(w http.ResponseWriter, r *http.Request) {
	client := CreateClientDTO{}
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		slog.Error("error decoding json input", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !client.Validate() {
		slog.Warn("data passed is invalid")
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
		slog.Error("error saving new client", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
