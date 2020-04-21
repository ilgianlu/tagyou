package main

import (
	"log"
	"os"

	mq "github.com/ilgianlu/tagyou/mqtt"
	dotenv "github.com/joho/godotenv"
)

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal(err)
		return
	}
	mq.StartMQTT(os.Getenv("LISTEN_PORT"))
}

func loadEnv() error {
	env := os.Getenv("TAGYOU_ENV")
	if "" == env {
		env = "default"
	}
	return dotenv.Load(".env." + env + ".local")
}
