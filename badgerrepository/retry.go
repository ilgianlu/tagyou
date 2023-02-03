package badgerrepository

import (
	"strconv"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

type RetryBadgerRepository struct {
	Db *badger.DB
}

func RetryKey(clientId string, packetIdentifier int, reasonCode uint8) []byte {
	k := []byte(clientId)
	k = append(k, []byte(strconv.Itoa(packetIdentifier))...)
	k = append(k, reasonCode)
	return k
}

func RetryPrefix(clientId string, packetIdentifier int) []byte {
	k := []byte(clientId)
	k = append(k, []byte(strconv.Itoa(packetIdentifier))...)
	return k
}

func RetryValue(retry model.Retry) ([]byte, error) {
	return model.GobEncode(retry)
}

func (r RetryBadgerRepository) SaveOne(retry model.Retry) error {
	return r.Db.Update(func(txn *badger.Txn) error {
		k := RetryKey(retry.ClientId, retry.PacketIdentifier, retry.ReasonCode)
		v, err := RetryValue(retry)
		if err != nil {
			return err
		}
		e := badger.NewEntry(k, v).WithTTL(time.Duration(int64(conf.RETRY_EXPIRATION)))
		return txn.SetEntry(e)
	})
}

func (r RetryBadgerRepository) Delete(retry model.Retry) error {
	return r.Db.Update(func(txn *badger.Txn) error {
		k := RetryKey(retry.ClientId, retry.PacketIdentifier, retry.ReasonCode)
		return txn.Delete(k)
	})
}

func (r RetryBadgerRepository) FirstByClientIdPacketIdentifier(clientId string, packetIdentifier int) (model.Retry, error) {
	var retry model.Retry
	r.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.ValidForPrefix(RetryPrefix(clientId, packetIdentifier)); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				ret, err := model.GobDecode[model.Retry](val)
				if err != nil {
					return err
				}
				retry = ret
				return nil
			})
		}
		return nil
	})
	return retry, nil
}

func (r RetryBadgerRepository) FirstByClientIdPacketIdentifierReasonCode(clientId string, packetIdentifier int, reasonCode uint8) (model.Retry, error) {
	key := RetryKey(clientId, packetIdentifier, reasonCode)
	ret := model.Retry{}

	err := r.Db.View(func(txn *badger.Txn) error {
		rItem, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		rItem.Value(func(val []byte) error {
			ret, err = model.GobDecode[model.Retry](val)
			return err
		})

		return err
	})

	return ret, err
}
