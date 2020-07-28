package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

type (
	// UserController represents the controller for operating on the User resource
	AuthController struct {
		db *gorm.DB
	}
)

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db}
}

func (uc AuthController) RegisterRoutes(r *httprouter.Router) {
	r.GET("/auths", uc.GetAuths)
	r.GET("/auths/:id", uc.GetAuth)
	r.POST("/auths", uc.CreateAuth)
	r.PUT("/auths/:id", uc.UpdateAuth)
	r.DELETE("/auths/:id", uc.RemoveAuth)
}

func (uc AuthController) GetAuths(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	auths := []model.Auth{}
	if err := uc.db.Find(&auths).Error; err != nil {
		log.Printf("error getting auth rows: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if res, err := json.Marshal(auths); err != nil {
		log.Printf("error marshaling auth rows: %s\n", err)
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

func (uc AuthController) GetAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	auth, err := uc.getOne(w, r, p)
	if err != nil {
		return
	}

	if res, err := json.Marshal(auth); err != nil {
		log.Printf("error marshaling auth row: %s\n", err)
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

// CreateUser creates a new user resource
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

// RemoveUser removes an existing user resource
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
		if gorm.IsRecordNotFoundError(err) {
			w.WriteHeader(http.StatusNoContent)
			return auth, fmt.Errorf("error getting auth row: %s\n", err)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			return auth, fmt.Errorf("error getting auth row: %s\n", err)
		}
	}

	return auth, nil
}
