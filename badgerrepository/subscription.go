package badgerrepository

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/model"
	"github.com/rs/zerolog/log"
)

type SubscriptionBadgerRepository struct {
	Db *badger.DB
}

func SubscriptionKey(clientId string, topic string, shareName string) []byte {
	k := []byte(topic)
	k = append(k, 0x00)
	k = append(k, []byte(clientId)...)
	k = append(k, []byte(shareName)...)
	return k
}

func SubscriptionValue(sub model.Subscription) ([]byte, error) {
	return model.GobEncode(sub)
}

func (s SubscriptionBadgerRepository) CreateOne(sub model.Subscription) error {
	return s.Db.Update(func(txn *badger.Txn) error {
		k := SubscriptionKey(sub.ClientId, sub.Topic, sub.ShareName)
		_, err := txn.Get(k)
		if err != badger.ErrKeyNotFound {
			return fmt.Errorf("subscription is already set for %s", sub.ClientId)
		}

		v, err := SubscriptionValue(sub)
		if err != nil {
			return err
		}
		return txn.Set(k, v)
	})
}

func (s SubscriptionBadgerRepository) FindToUnsubscribe(shareName string, topic string, clientId string) (model.Subscription, error) {
	key := SubscriptionKey(clientId, topic, shareName)
	sub := model.Subscription{}

	err := s.Db.View(func(txn *badger.Txn) error {
		sItem, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		sItem.Value(func(val []byte) error {
			sub, err = model.GobDecode[model.Subscription](val)
			return err
		})

		return err
	})

	return sub, err
}

func (s SubscriptionBadgerRepository) FindSubscriptions(topics []string, shared bool) []model.Subscription {
	subscriptions := []model.Subscription{}
	for _, tpc := range topics {
		subs := s.findBySubTopic(tpc)
		subscriptions = append(subscriptions, subs...)
	}
	result := []model.Subscription{}
	for _, s := range subscriptions {
		if s.Shared == shared {
			result = append(result, s)
		}
	}
	return result
}

func (s SubscriptionBadgerRepository) findBySubTopic(topic string) []model.Subscription {
	subscriptions := []model.Subscription{}
	prefix := []byte(topic)
	prefix = append(prefix, 0x00)
	s.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				sub, err := model.GobDecode[model.Subscription](val)
				if err != nil {
					log.Err(err).Msg("error during gob decode of subscription")
					return err
				}
				subscriptions = append(subscriptions, sub)
				return nil
			})
		}
		return nil
	})
	return subscriptions
}

func (s SubscriptionBadgerRepository) DeleteByClientIdTopicShareName(clientId string, topic string, shareName string) error {
	return s.Db.Update(func(txn *badger.Txn) error {
		k := SubscriptionKey(clientId, topic, shareName)
		return txn.Delete(k)
	})
}

func (s SubscriptionBadgerRepository) GetAll() []model.Subscription {
	subscriptions := []model.Subscription{}
	s.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				sub, err := model.GobDecode[model.Subscription](val)
				if err != nil {
					return err
				}
				subscriptions = append(subscriptions, sub)
				return nil
			})
		}
		return nil
	})
	return subscriptions
}
