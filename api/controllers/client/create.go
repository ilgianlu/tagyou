package client

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

func (uc ClientController) CreateClient(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	client := model.Client{}
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		log.Error().Err(err).Msg("error decoding json input")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !client.Validate() || !client.ValidPassword() {
		log.Error().Msg("data passed is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := client.SetPassword(client.InputPassword); err != nil {
		log.Error().Err(err).Msg("error encoding password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := persistence.ClientRepository.Create(client); err != nil {
		log.Error().Err(err).Msg("error saving new client")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	client.InputPassword = ""
	client.InputPasswordConfirm = ""
	if res, err := json.Marshal(client); err != nil {
		log.Error().Err(err).Msg("error marshaling new client")
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
		log.Debug().Msgf("Wrote %d bytes json result\n", numBytes)
	}
}
