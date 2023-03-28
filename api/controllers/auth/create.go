package auth

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/jwt"
	"github.com/ilgianlu/tagyou/password"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

type CreateTokenDTO struct {
	Username      string
	InputPassword string
}

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (uc AuthController) CreateToken(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	user := CreateTokenDTO{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error().Err(err).Msg("error decoding json input")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	usr, err := persistence.UserRepository.GetByUsername(user.Username)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := password.CheckPassword(usr.Password, []byte(user.InputPassword)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := jwt.CreateToken(
		conf.API_TOKEN_SIGNING_KEY,
		conf.API_TOKEN_ISSUER,
		conf.API_TOKEN_HOURS_DURATION,
		usr.ID,
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	response := CreateTokenResponse{Token: token}

	if res, err := json.Marshal(response); err != nil {
		log.Printf("error marshaling token response: %s\n", err)
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
