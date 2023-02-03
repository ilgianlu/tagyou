package badgerrepository

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/topic"
)

type RetainBadgerRepository struct {
	Db *badger.DB
}

func RetainKey(topic string) []byte {
	k := []byte(topic)
	return k
}

func RetainValue(retain model.Retain) ([]byte, error) {
	return model.GobEncode(retain)
}

func (r RetainBadgerRepository) FindRetains(subscribedTopic string) []model.Retain {
	var retains []model.Retain
	r.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			if topic.Match(string(item.Key()), subscribedTopic) {
				item.Value(func(val []byte) error {
					ret, err := model.GobDecode[model.Retain](val)
					if err != nil {
						return err
					}
					retains = append(retains, ret)
					return nil
				})
			}
		}
		return nil
	})
	return retains
}

func (r RetainBadgerRepository) Create(retain model.Retain) error {
	return r.Db.Update(func(txn *badger.Txn) error {
		k := RetainKey(retain.Topic)
		v, err := RetainValue(retain)
		if err != nil {
			return err
		}
		return txn.Set(k, v)
	})
}

func (r RetainBadgerRepository) Delete(retain model.Retain) error {
	return r.Db.Update(func(txn *badger.Txn) error {
		k := RetainKey(retain.Topic)
		return txn.Delete(k)
	})
}