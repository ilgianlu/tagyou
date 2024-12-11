package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

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

	conf.Loader()

	log.Init()

	p := persistence.SqlPersistence{
		DbFile:       conf.DB_PATH + "/" + conf.DB_NAME,
		InitDatabase: conf.INIT_DB,
	}
	p.Init(conf.CLEAN_EXPIRED_SESSIONS, conf.CLEAN_EXPIRED_RETRIES, conf.INIT_ADMIN_PASSWORD)

	router, debugFile := selectRouter(conf.ROUTER_MODE)
	if debugFile != nil {
		defer debugFile.Close()
	}

	go api.StartApi(conf.API_PORT)
	go mqtt.StartWebSocket(conf.WS_PORT, router)
	go mqtt.StartMQTT(conf.LISTEN_PORT, router)

	<-c
}

func selectRouter(mode string) (routers.Router, *os.File) {
	slog.Info("starting with router", "mode", strings.ToUpper(mode))
	switch mode {
	case conf.ROUTER_MODE_DEBUG:
		file, err := os.Create(conf.DEBUG_FILE)
		if err != nil {
			slog.Error("Could write debug file", "err", err)
			panic(1)
		}
		return routers.NewDebug(file, conf.DEBUG_CLIENTS), file
	case conf.ROUTER_MODE_SIMPLE:
		return routers.NewSimple(), nil
	default:
		return routers.NewStandard(), nil
	}
}
