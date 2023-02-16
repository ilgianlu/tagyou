package api

import (
	"net/http"

	"github.com/rs/zerolog/log"

	AuthController "github.com/ilgianlu/tagyou/api/controllers/auth"
	SessionController "github.com/ilgianlu/tagyou/api/controllers/session"
	SubscriptionController "github.com/ilgianlu/tagyou/api/controllers/subscription"
	"github.com/julienschmidt/httprouter"
)

func StartApi(httpPort string) {
	r := httprouter.New()
	uc := AuthController.New()
	uc.RegisterRoutes(r)
	sc := SessionController.New()
	sc.RegisterRoutes(r)
	subc := SubscriptionController.New()
	subc.RegisterRoutes(r)

	log.Info().Msgf("[API] http listening on %s", httpPort)
	if err := http.ListenAndServe(httpPort, r); err != nil {
		log.Fatal().Err(err).Msg("[API] http listener broken")
	}
}
