package auth

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

func (uc AuthController) CreateAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	auth := model.Auth{}
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		log.Error().Err(err).Msg("error decoding json input")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !auth.Validate() || !auth.ValidPassword() {
		log.Error().Msg("data passed is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := auth.SetPassword(auth.InputPassword); err != nil {
		log.Error().Err(err).Msg("error encoding password")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := persistence.AuthRepository.Create(auth); err != nil {
		log.Error().Err(err).Msg("error saving new auth")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	auth.InputPassword = ""
	auth.InputPasswordConfirm = ""
	if res, err := json.Marshal(auth); err != nil {
		log.Error().Err(err).Msg("error marshaling new auth")
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
