package subscription

import (
	"github.com/ilgianlu/tagyou/api/controllers/middlewares"
	"github.com/julienschmidt/httprouter"
)

const resourceName string = "/subscriptions"

type SubscriptionController struct {
}

func New() *SubscriptionController {
	return &SubscriptionController{}
}

func (sc SubscriptionController) RegisterRoutes(r *httprouter.Router) {
	r.GET(resourceName, middlewares.Authenticated(sc.GetSubscriptions))
}
