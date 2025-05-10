package insights

import (
	"fmt"
	"net/http"

	"github.com/ilgianlu/tagyou/conf"
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

func blueprintFilepath(model string) string {
	return fmt.Sprintf(conf.AI_BLUEPRINT_PATH + "/%s.csv", model)
}

func debugDataFilepath(clientID string) string {
	return fmt.Sprintf(conf.DEBUG_DATA_PATH + "/%s.dump", clientID)
}
