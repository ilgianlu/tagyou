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
	"github.com/jinzhu/gorm"
	"github.com/julienschmidt/httprouter"
)

func StartApi(httpPort string) {
	db, err := gorm.Open("sqlite3", os.Getenv("DB_PATH")+os.Getenv("DB_NAME"))
	if err != nil {
		log.Fatalf("[API] failed to connect database %s", err)
	}
	log.Println("[API] db connected !")
	// db.LogMode(true)
	defer db.Close()

	clientOptions := mqtt.NewClientOptions().
		SetClientID("api").
		AddBroker(os.Getenv("LISTEN_PORT")).
		SetConnectionLostHandler(connLostHandler).
		SetConnectTimeout(1 * time.Second).
		SetOnConnectHandler(onConnectHandler)

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

func mqttConnect(c mqtt.Client) {
	i := 0
	success := false
	for !success {
		token := c.Connect()
		token.WaitTimeout(5 * time.Second)
		if token.Wait() && token.Error() != nil {
			log.Println(token.Error())
		} else {
			success = true
		}
		if i == 3 {
			panic(token.Error())
		}
		i = i + 1
	}
}

func connLostHandler(c mqtt.Client, err error) {
	log.Printf("Connection lost, reason: %v\n", err)
	//Perform additional action...
}

func onConnectHandler(c mqtt.Client) {
	log.Println("Client Connected")
	//Perform additional action...
}
