package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/users"

type UserController struct {
}

func New() *UserController {
	return &UserController{}
}

func (uc UserController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, uc.GetUsers)
	r.POST(resourceName, uc.CreateUser)
}

func (uc UserController) getOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) (model.User, error) {
	id := p.ByName("id")

	idNum, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return model.User{}, fmt.Errorf("invalid user id")
	}

	user, err := persistence.UserRepository.GetById(uint(idNum))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return user, fmt.Errorf("error getting user row: %s", err)
	}

	return user, nil
}
