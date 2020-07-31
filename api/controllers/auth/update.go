package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (uc AuthController) UpdateAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	auth, err := uc.getOne(w, r, p)
	if err != nil {
		return
	}

	update := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		log.Printf("error decoding json input: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tx := uc.db.Begin()
	tx.Model(&auth).Omit("password").Updates(update)

	if auth.Validate() {
		tx.Commit()
		w.WriteHeader(http.StatusOK)
	} else {
		tx.Rollback()
		w.WriteHeader(http.StatusBadRequest)
	}
}
