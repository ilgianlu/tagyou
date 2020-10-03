package api

import (
	"log"
	"net/http"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	AuthController "github.com/ilgianlu/tagyou/api/controllers/auth"
	MessageController "github.com/ilgianlu/tagyou/api/controllers/message"
	SessionController "github.com/ilgianlu/tagyou/api/controllers/session"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func StartApi(httpPort string) {
	db, err := gorm.Open(sqlite.Open(os.Getenv("DB_PATH")+os.Getenv("DB_NAME")), &gorm.Config{})
	if err != nil {
		log.Fatalf("[API] failed to connect database %s", err)
	}
	log.Println("[API] db connected !")
	// db.LogMode(true)
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

	log.Printf("[API] http listening on %s", httpPort)
	if err := http.ListenAndServe(httpPort, r); err != nil {
		log.Panic(err)
	}
}

func closeDb(db *gorm.DB) {
	sql, err := db.DB()
	if err != nil {
		log.Println("could not close DB", err)
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
			log.Printf("[API] mqtt connect error %s\n", token.Error())
		} else {
			success = true
		}
		if i == 3 {
			log.Printf("[API] panicking after too many connect errors %s\n", token.Error())
			panic(token.Error())
		}
		i = i + 1
	}
}

func connLostHandler(c mqtt.Client, err error) {
	log.Printf("[API] MQTT Connection lost, reason: %v\n", err)
	//Perform additional action...
}

func onConnectHandler(c mqtt.Client) {
	log.Println("[API] MQTT Client Connected")
	//Perform additional action...
}
