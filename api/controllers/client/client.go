package client

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ilgianlu/tagyou/api/controllers/middlewares"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
)

type ClientController struct {
}

func NewController() *ClientController {
	return &ClientController{}
}

func (uc ClientController) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /clients", middlewares.Authenticated(uc.GetClients))
	r.HandleFunc("GET /clients/{id}", middlewares.Authenticated(uc.GetClient))
	r.HandleFunc("POST /clients", middlewares.Authenticated(uc.CreateClient))
	r.HandleFunc("DELETE /clients/{id}", middlewares.Authenticated(uc.DeleteClient))
}

func (uc ClientController) getOne(w http.ResponseWriter, r *http.Request) (model.Client, error) {
	id := r.PathValue("id")

	idNum, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return model.Client{}, fmt.Errorf("invalid user id")
	}

	client, err := persistence.ClientRepository.GetById(int64(idNum))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return client, fmt.Errorf("error getting client row: %s", err)
	}

	return client, nil
}
