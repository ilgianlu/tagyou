package api

import (
	"log/slog"
	"net/http"

	"github.com/ilgianlu/tagyou/api/controllers/auth"
	"github.com/ilgianlu/tagyou/api/controllers/client"
	"github.com/ilgianlu/tagyou/api/controllers/session"
	"github.com/ilgianlu/tagyou/api/controllers/subscription"
	"github.com/ilgianlu/tagyou/api/controllers/user"
)

func StartApi(httpPort string) {
	mux := http.NewServeMux()
	authController := auth.NewController()
	authController.RegisterRoutes(mux)
	clientController := client.NewController()
	clientController.RegisterRoutes(mux)
	sc := session.NewController()
	sc.RegisterRoutes(mux)
	subc := subscription.NewController()
	subc.RegisterRoutes(mux)
	usrc := user.NewController()
	usrc.RegisterRoutes(mux)

	slog.Info("[API] http start listening", "port", httpPort)
	if err := http.ListenAndServe(httpPort, mux); err != nil {
		slog.Error("[API] http listener broken", "err", err)
		panic(1)
	}
}
