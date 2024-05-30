package subscription

import (
	"net/http"

	"github.com/ilgianlu/tagyou/api/controllers/middlewares"
)

type SubscriptionController struct {
}

func NewController() *SubscriptionController {
	return &SubscriptionController{}
}

func (sc SubscriptionController) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("GET /subscriptions", middlewares.Authenticated(sc.GetSubscriptions))
}
