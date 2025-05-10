package insights

import (
	"net/http"

	"github.com/ilgianlu/tagyou/api/controllers/middlewares"
)

type InsightsController struct {
}

func NewController() *InsightsController {
	return &InsightsController{}
}

func (dc InsightsController) RegisterRoutes(r *http.ServeMux) {
	r.HandleFunc("POST /insights/report", middlewares.Authenticated(dc.Report))
	r.HandleFunc("POST /insights/blueprint", middlewares.Authenticated(dc.Blueprint))
}

