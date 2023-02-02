package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/auths"

type AuthController struct {
}

func New() *AuthController {
	return &AuthController{}
}

func (uc AuthController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, uc.GetAuths)
	r.GET(resourceName+"/:id", uc.GetAuth)
	r.POST(resourceName, uc.CreateAuth)
	r.DELETE(resourceName+"/:id", uc.RemoveAuth)
}

func (uc AuthController) getOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) (model.Auth, error) {
	id := p.ByName("id")

	idParts := strings.Split(id, "-")
	if len(idParts) != 2 {
		w.WriteHeader(http.StatusBadRequest)
		return model.Auth{}, fmt.Errorf("invalid auth id")
	}

	auth, err := persistence.AuthRepository.GetByClientIdUsername(idParts[0], idParts[1])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return auth, fmt.Errorf("error getting auth row: %s", err)
	}

	return auth, nil
}
