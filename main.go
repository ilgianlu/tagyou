package main

import (
	"fmt"
	"os"

	mq "github.com/ilgianlu/tagyou/mqtt"
	dotenv "github.com/joho/godotenv"
)

func main() {
	// load env vars
	berr := dotenv.Load()
	if berr != nil {
		fmt.Println("error loading env", berr)
	}

	mq.StartMQTT(os.Getenv("LISTEN_PORT"))
}
