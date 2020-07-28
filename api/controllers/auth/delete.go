package auth

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (uc AuthController) RemoveAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	auth, err := uc.getOne(w, r, p)
	if err != nil {
		return
	}
	if err := uc.db.Delete(&auth).Error; err != nil {
		log.Printf("error deleting auth row: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
