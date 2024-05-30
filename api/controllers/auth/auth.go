package auth

import (
	"net/http"
)

type AuthController struct {
}

func NewController() *AuthController {
	return &AuthController{}
}

func (ac AuthController) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /auth", ac.CreateToken)
}
