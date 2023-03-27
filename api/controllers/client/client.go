package client

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/clients"

type ClientController struct {
}

func New() *ClientController {
	return &ClientController{}
}

func (uc ClientController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, uc.GetClients)
	r.GET(resourceName+"/:id", uc.GetClient)
	r.POST(resourceName, uc.CreateClient)
	r.DELETE(resourceName+"/:id", uc.DeleteClient)
}

func (uc ClientController) getOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) (model.Client, error) {
	id := p.ByName("id")

	idNum, err := strconv.Atoi(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return model.Client{}, fmt.Errorf("invalid user id")
	}

	client, err := persistence.ClientRepository.GetById(uint(idNum))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return client, fmt.Errorf("error getting client row: %s", err)
	}

	return client, nil
}
