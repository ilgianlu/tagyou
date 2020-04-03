package mqtt

import (
	"fmt"

	bolt "go.etcd.io/bbolt"
)

const SUBSCRIPTION_BUCKET = "subscriptions"

type Subscription struct {
	topic    string
	clientId string
}

func (s Subscription) persist(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		subscriptions := tx.Bucket([]byte(SUBSCRIPTION_BUCKET))
		serr := subscriptions.Put([]byte(s.topic), []byte(s.clientId))
		if serr != nil {
			fmt.Println(serr)
			return nil
		}
		return nil
	})
}

func (s Subscription) remove(db *bolt.DB) error {
	return db.Update(func(tx *bolt.Tx) error {
		subscriptions := tx.Bucket([]byte(SUBSCRIPTION_BUCKET))
		serr := subscriptions.Delete([]byte(s.topic))
		if serr != nil {
			fmt.Println(serr)
			return nil
		}
		return nil
	})
}

func findSubs(db *bolt.DB, topic string) []string {
	clientIds := make([]string, 0)
	ss := make(chan string)
	go db.View(func(tx *bolt.Tx) error {
		subscriptions := tx.Bucket([]byte(SUBSCRIPTION_BUCKET))
		v := subscriptions.Get([]byte(topic))
		if v != nil {
			ss <- string(v)
		}
		ss <- ""
		return nil
	})
	return append(clientIds, <-ss)
}
