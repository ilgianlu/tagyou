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

func (c ConnStatus) persist(db *bolt.DB) error {
	return db.Batch(func(tx *bolt.Tx) error {
		clients, err0 := tx.CreateBucketIfNotExists([]byte(CLIENTS_BUCKET))
		if err0 != nil {
			fmt.Println(err0)
			return err0
		}
		client, err1 := clients.CreateBucketIfNotExists([]byte(c.clientId))
		if err1 != nil {
			fmt.Println(err1)
			return err1
		}
		client.Put([]byte(CONNECT_TIME), []byte(c.connectTime.Format(time.UnixDate)))
		client.Put([]byte(CLIENTID), []byte(c.clientId))
		client.Put([]byte(PROTOCOL_VERSION), []byte{c.protocolVersion})
		client.Put([]byte(CONNECT_FLAGS), []byte{c.connectFlags})
		client.Put([]byte(KEEP_ALIVE), []byte(c.keepAlive))
		return nil
	})
}
