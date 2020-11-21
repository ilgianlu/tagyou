package api

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog/log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	AuthController "github.com/ilgianlu/tagyou/api/controllers/auth"
	MessageController "github.com/ilgianlu/tagyou/api/controllers/message"
	SessionController "github.com/ilgianlu/tagyou/api/controllers/session"
	SubscriptionController "github.com/ilgianlu/tagyou/api/controllers/subscription"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func StartApi(httpPort string) {
	db, err := openDb()
	if err != nil {
		log.Fatal().Err(err).Msgf("[API] failed to connect database %s", os.Getenv("DB_PATH")+os.Getenv("DB_NAME"))
	}
	log.Info().Msg("[API] db connected !")
	defer closeDb(db)

	clientOptions := mqtt.NewClientOptions().
		SetClientID("api").
		AddBroker(os.Getenv("LISTEN_PORT")).
		SetConnectionLostHandler(connLostHandler).
		SetConnectTimeout(1 * time.Second).
		SetOnConnectHandler(onConnectHandler).
		SetKeepAlive(time.Duration(conf.DEFAULT_KEEPALIVE) * time.Second)

	// mqtt.DEBUG = log.New(os.Stderr, "DEBUG    ", log.Ltime)
	c := mqtt.NewClient(clientOptions)
	go mqttConnect(c)

	r := httprouter.New()
	uc := AuthController.New(db)
	uc.RegisterRoutes(r)
	sc := SessionController.New(db)
	sc.RegisterRoutes(r)
	mc := MessageController.New(c)
	mc.RegisterRoutes(r)
	subc := SubscriptionController.New(db)
	subc.RegisterRoutes(r)

	log.Info().Msgf("[API] http listening on %s", httpPort)
	if err := http.ListenAndServe(httpPort, r); err != nil {
		log.Fatal().Err(err).Msg("[API] http listener broken")
	}
}

func openDb() (*gorm.DB, error) {
	logLevel := logger.Silent
	if os.Getenv("DEBUG") != "" {
		logLevel = logger.Info
	}
	return gorm.Open(sqlite.Open(os.Getenv("DB_PATH")+os.Getenv("DB_NAME")), &gorm.Config{
		Logger: logger.New(
			&log.Logger,
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      logLevel,
				Colorful:      true,
			},
		),
	})
}

func closeDb(db *gorm.DB) {
	sql, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("[API] could not close DB")
		return
	}
	sql.Close()
}

func mqttConnect(c mqtt.Client) {
	time.Sleep(5 * time.Second)
	i := 0
	success := false
	for !success {
		token := c.Connect()
		token.WaitTimeout(5 * time.Second)
		if token.Wait() && token.Error() != nil {
			log.Error().Err(token.Error()).Msg("[API] mqtt connect error")
		} else {
			success = true
		}
		if i == 3 {
			log.Fatal().Err(token.Error()).Msg("[API] panicking after too many connect errors")
			panic(token.Error())
		}
		i = i + 1
	}
}

func connLostHandler(c mqtt.Client, err error) {
	log.Debug().Err(err).Msg("[API] MQTT Connection lost")
	//Perform additional action...
}

func onConnectHandler(c mqtt.Client) {
	log.Debug().Msg("[API] MQTT Client Connected")
	//Perform additional action...
}
