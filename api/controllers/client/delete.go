package client

import (
	"log/slog"
	"net/http"

	"github.com/ilgianlu/tagyou/persistence"

	"github.com/julienschmidt/httprouter"
)

func (uc ClientController) DeleteClient(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	client, err := uc.getOne(w, r, p)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := persistence.ClientRepository.DeleteById(client.ID); err != nil {
		slog.Error("error deleting auth row", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
