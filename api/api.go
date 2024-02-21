package api

import (
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"

	AuthController "github.com/ilgianlu/tagyou/api/controllers/auth"
	ClientController "github.com/ilgianlu/tagyou/api/controllers/client"
	SessionController "github.com/ilgianlu/tagyou/api/controllers/session"
	SubscriptionController "github.com/ilgianlu/tagyou/api/controllers/subscription"
	UserController "github.com/ilgianlu/tagyou/api/controllers/user"
)

func StartApi(httpPort string) {
	r := httprouter.New()
	ac := AuthController.New()
	ac.RegisterRoutes(r)
	uc := ClientController.New()
	uc.RegisterRoutes(r)
	sc := SessionController.New()
	sc.RegisterRoutes(r)
	subc := SubscriptionController.New()
	subc.RegisterRoutes(r)
	usrc := UserController.New()
	usrc.RegisterRoutes(r)

	slog.Info("[API] http start listening", "port", httpPort)
	if err := http.ListenAndServe(httpPort, r); err != nil {
		slog.Error("[API] http listener broken", "err", err)
		panic(1)
	}
}
