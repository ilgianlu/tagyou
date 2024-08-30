package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	dotenv "github.com/joho/godotenv"

	"github.com/ilgianlu/tagyou/api"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/log"
	"github.com/ilgianlu/tagyou/mqtt"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	err := loadEnv()
	if err != nil {
		slog.Error("Could not load env", "err", err)
		panic(1)
	}
	conf.Loader()

	log.Init()

	p := persistence.SqlPersistence{}
	p.Init()

	var router routers.Router
	if conf.ROUTER_MODE == "debug" {
		file, err := os.Create(conf.DEBUG_FILE)
		if err != nil {
			slog.Error("Could write debug file", "err", err)
			panic(1)
		}
		defer file.Close()
		router = routers.NewDebug(file)
	} else if conf.ROUTER_MODE == "simple" {
		router = routers.NewSimple()
	} else {
		router = routers.NewStandard()
	}

	go api.StartApi(os.Getenv("API_PORT"))
	go mqtt.StartWebSocket(os.Getenv("WS_PORT"), router)
	go mqtt.StartMQTT(os.Getenv("LISTEN_PORT"), router)

	<-c
}

func loadEnv() error {
	env := os.Getenv("TAGYOU_ENV")
	if env == "" {
		env = "default"
	}
	return dotenv.Load(".env." + env + ".local")
}
