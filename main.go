package main

import (
	"fmt"
	"os"

	"github.com/ilgianlu/tagyou/mqtt"
	dotenv "github.com/joho/godotenv"
	bolt "go.etcd.io/bbolt"
)

func main() {
	// load env vars
	berr := dotenv.Load()
	if berr != nil {
		fmt.Println("error loading env", berr)
	}

	// open k/v store
	db, derr := bolt.Open(os.Getenv("DB_PATH"), 0666, nil)
	if derr != nil {
		fmt.Println("error opening bbolt", derr)
	}
	defer db.Close()

	mq := mqtt.New(db)
	mq.Start(os.Getenv("LISTEN_PORT"))
}
