package mqtt

import (
	"fmt"
	"time"

	bolt "go.etcd.io/bbolt"
)

const CLIENTS_BUCKET = "clients"
const CONNECT_TIME = "connectTime"
const CLIENTID = "clientId"
const PROTOCOL_VERSION = "protocolVersion"
const CONNECT_FLAGS = "connectFlags"
const KEEP_ALIVE = "keepAlive"

type ConnStatus struct {
	connectTime     time.Time
	clientId        string
	protocolVersion uint8
	connectFlags    uint8
	keepAlive       []byte
}

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

func newClient(db *bolt.DB, connStatus *ConnStatus) {
	c := make(chan string)
	go db.Batch(func(tx *bolt.Tx) error {
		clients := tx.Bucket([]byte(CLIENTS_BUCKET))
		client, err := clients.CreateBucketIfNotExists([]byte(connStatus.clientId))
		if err != nil {
			c <- connStatus.clientId
			return nil
		}
		client.Put([]byte(CONNECT_TIME), []byte(connStatus.connectTime.Format(time.UnixDate)))
		client.Put([]byte(CLIENTID), []byte(connStatus.clientId))
		client.Put([]byte(PROTOCOL_VERSION), []byte{connStatus.protocolVersion})
		client.Put([]byte(CONNECT_FLAGS), []byte{connStatus.connectFlags})
		client.Put([]byte(KEEP_ALIVE), []byte(connStatus.keepAlive))
		if err != nil {
			c <- connStatus.clientId
			return nil
		}
		c <- "0"
		return nil
	})
	fmt.Println(<-c)
}
