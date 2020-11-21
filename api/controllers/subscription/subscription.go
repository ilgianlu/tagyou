package subscription

import (
	"github.com/julienschmidt/httprouter"
	"gorm.io/gorm"
)

const resourceName string = "/subscriptions"

type SubscriptionController struct {
	db *gorm.DB
}

func New(db *gorm.DB) *SubscriptionController {
	return &SubscriptionController{db}
}

func (sc SubscriptionController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, sc.GetSubscriptions)
}
