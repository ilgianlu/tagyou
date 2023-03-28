package auth

import (
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/auth"

type AuthController struct {
}

func New() *AuthController {
	return &AuthController{}
}

func (uc AuthController) RegisterRoutes(r *httprouter.Router) {
	r.POST(resourceName, uc.CreateToken)
}
