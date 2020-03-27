package mqtt

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

const CLIENTS_BUCKET = "clients"

func initBucket(db *bolt.DB) error {
	uerr := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(CLIENTS_BUCKET))
		if err != nil {
			return err
		}
		return nil
	})
	return uerr
}

func newClient(db *bolt.DB, clientId string) {
	c := make(chan string)
	go db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CLIENTS_BUCKET))
		err := b.Put([]byte(clientId), []byte(time.Now().Format(time.UnixDate)))
		if err != nil {
			c <- clientId
			return nil
		}
		c <- "0"
		return nil
	})
	fmt.Println(<-c)
}

func clientOk(db *bolt.DB, clientId string) {
	c := make(chan string)
	go db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(CLIENTS_BUCKET))
		v := b.Get([]byte(clientId))
		if v != nil {
			c <- clientId
			return nil
		}
		c <- "0"
		return nil
	})
	fmt.Println(<-c)
}
