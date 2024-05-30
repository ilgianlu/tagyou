package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/jwt"
	"github.com/ilgianlu/tagyou/password"
	"github.com/ilgianlu/tagyou/persistence"
)

type CreateTokenDTO struct {
	Username      string
	InputPassword string
}

type CreateTokenResponse struct {
	Token string `json:"token"`
}

func (ac AuthController) CreateToken(w http.ResponseWriter, r *http.Request) {
	user := CreateTokenDTO{}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		slog.Error("error decoding json input", "err", err)
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
		slog.Error("error marshaling token response", "err", err)
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
