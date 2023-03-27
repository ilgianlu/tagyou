package user

import (
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/users"

type UserController struct {
}

func New() *UserController {
	return &UserController{}
}

func (uc UserController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, uc.GetUsers)
	// r.GET(resourceName+"/:id", uc.GetClient)
	// r.POST(resourceName, uc.CreateClient)
	// r.DELETE(resourceName+"/:id", uc.RemoveClient)
}
