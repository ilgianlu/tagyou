package auth

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ilgianlu/tagyou/model"
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

const resourceName string = "/auths"

type AuthController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *AuthController {
	return &AuthController{db}
}

func (uc AuthController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, uc.GetAuths)
	r.GET(resourceName+"/:id", uc.GetAuth)
	r.POST(resourceName, uc.CreateAuth)
	r.PUT(resourceName+"/:id", uc.UpdateAuth)
	r.DELETE(resourceName+"/:id", uc.RemoveAuth)
}

func (uc AuthController) getOne(w http.ResponseWriter, r *http.Request, p httprouter.Params) (model.Auth, error) {
	auth := model.Auth{}

	id := p.ByName("id")
	authId, err := strconv.Atoi(id)
	if err != nil {
		log.Printf("passing bad id: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return auth, fmt.Errorf("passing bad id: %s\n", err)
	}

	if err := uc.db.Where("id = ?", authId).First(&auth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNoContent)
			return auth, fmt.Errorf("error getting auth row: %s\n", err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return auth, fmt.Errorf("error getting auth row: %s\n", err)
		}
	}

	return auth, nil
}
