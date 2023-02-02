package badgerrepository

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

const RETAIN_PREFIX = "retain@"

type RetainBadgerRepository struct {
	Db *badger.DB
}

func RetainKey(topic string) []byte {
	k := []byte(RETAIN_PREFIX)
	k = append(k, []byte(topic)...)
	return k
}

func RetainValue(retain model.Retain) ([]byte, error) {
	return model.GobEncode(retain)
}

func (r RetainBadgerRepository) FindRetains(subscribedTopic string) []model.Retain {
	trimmedTopic := trimWildcard(subscribedTopic)
	var retains []model.Retain

	r.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				a, err := model.GobDecode[model.Retain](v)
				if err != nil {
					return err
				}
				retains = append(retains, a)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return retains
}

func (r RetainBadgerRepository) Create(retain model.Retain) error {
	return r.Db.Create(&retain).Error
}

func (r RetainBadgerRepository) Delete(retain model.Retain) error {
	return r.Db.Delete(&retain).Error
}

func trimWildcard(topic string) string {
	lci := len(topic) - 1
	lc := topic[lci]
	if string(lc) == conf.WILDCARD_MULTI_LEVEL {
		topic = topic[:lci]
	}
	return topic
}
