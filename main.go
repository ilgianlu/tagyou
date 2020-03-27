package main

import (
	"fmt"
	"net"
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

	// start tcp socket
	ln, err := net.Listen("tcp", os.Getenv("LISTEN_PORT"))
	if err != nil {
		// handle error
		fmt.Println("error", err)
	}
	fmt.Println("listen", os.Getenv("LISTEN_PORT"))

	mq := mqtt.New(db)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Println("error ", err)
		}
		go mq.HandleConnection(conn)
	}
}
