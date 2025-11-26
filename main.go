package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ilgianlu/tagyou/api"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/log"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/mqtt"
	"github.com/ilgianlu/tagyou/persistence"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	conf.Loader()

	logLevel := log.Init(os.Getenv("DEBUG") != "")
	slog.Warn("Configuration loaded, Logging started, Tagyou starting up", "log_level", logLevel)

	p := persistence.SqlPersistence{
		DbFile:       conf.DB_PATH + "/" + conf.DB_NAME,
		InitDatabase: conf.INIT_DB,
	}
	_, err := p.Init(conf.CLEAN_EXPIRED_SESSIONS, conf.CLEAN_EXPIRED_RETRIES, conf.INIT_ADMIN_PASSWORD)
	if err != nil {
		panic(1)
	}

	connections := model.SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	go api.StartApi(conf.API_PORT)
	go mqtt.StartWebSocket(conf.WS_PORT, &connections)
	go mqtt.StartMQTT(conf.LISTEN_PORT, &connections)

	<-ctx.Done()
	slog.Warn("Going away...")
	p.Close()
}
