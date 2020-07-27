package controllers

import (
	"fmt"
	"net/http"

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

func (uc AuthController) GetAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("id")
	fmt.Println("asking for", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "")
}

// CreateUser creates a new user resource
func (uc AuthController) CreateAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
}

// RemoveUser removes an existing user resource
func (uc AuthController) RemoveAuth(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
}
