package auth

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ilgianlu/tagyou/model"
	"github.com/julienschmidt/httprouter"
)

func (uc AuthController) CreateAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	auth := model.Auth{}
	if err := json.NewDecoder(r.Body).Decode(&auth); err != nil {
		log.Printf("error decoding json input: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !auth.Validate() || !auth.ValidPassword() {
		log.Println("data passed is invalid")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := auth.SetPassword(auth.InputPassword); err != nil {
		log.Printf("error encoding password: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := uc.db.Save(&auth).Error; err != nil {
		log.Printf("error saving new auth: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	auth.InputPassword = ""
	auth.InputPasswordConfirm = ""
	if res, err := json.Marshal(auth); err != nil {
		log.Printf("error marshaling new auth: %s\n", err)
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
