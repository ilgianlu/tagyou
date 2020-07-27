package main

import (
	"log"
	"os"

	"github.com/ilgianlu/tagyou/api"
	mq "github.com/ilgianlu/tagyou/mqtt"
	dotenv "github.com/joho/godotenv"
)

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
		return
	}

	go api.StartApi(os.Getenv("API_PORT"))
	mq.StartMQTT(os.Getenv("LISTEN_PORT"))
}

func loadEnv() error {
	env := os.Getenv("TAGYOU_ENV")
	if "" == env {
		env = "default"
	}
	return dotenv.Load(".env." + env + ".local")
}
