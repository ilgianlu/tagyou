package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ilgianlu/tagyou/api/controllers/middlewares"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
)

type UserController struct {
}

func NewController() *UserController {
	return &UserController{}
}

func (uc UserController) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /users", middlewares.Authenticated(uc.GetUsers))
	r.HandleFunc("POST /users", middlewares.Authenticated(uc.CreateUser))
	r.HandleFunc("DELETE /users/{id}", middlewares.Authenticated(uc.DeleteUser))
}

func (uc UserController) getOne(w http.ResponseWriter, r *http.Request) (model.User, error) {
	id := r.PathValue("id")

	idNum, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return model.User{}, fmt.Errorf("invalid user id")
	}

	user, err := persistence.UserRepository.GetById(int64(idNum))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return user, fmt.Errorf("error getting user row: %s", err)
	}

	return user, nil
}
