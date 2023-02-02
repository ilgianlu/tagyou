package badgerrepository

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/model"
)

const AUTH_PREFIX = "auth@"

func AuthKey(clientId string, username string) []byte {
	k := []byte(AUTH_PREFIX)
	k = append(k, []byte(clientId)...)
	k = append(k, []byte(username)...)
	return k
}

func AuthValue(auth model.Auth) ([]byte, error) {
	return model.GobEncode(auth)
}

type AuthBadgerRepository struct {
	Db *badger.DB
}

func (ar AuthBadgerRepository) GetByClientIdUsername(clientId string, username string) (model.Auth, error) {
	key := AuthKey(clientId, username)
	a := model.Auth{}

	err := ar.Db.View(func(txn *badger.Txn) error {
		aItem, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		aItem.Value(func(val []byte) error {
			a, err = model.GobDecode[model.Auth](val)
			return err
		})

		return err
	})

	return a, err
}

func (ar AuthBadgerRepository) Create(auth model.Auth) error {
	return ar.Db.Update(func(txn *badger.Txn) error {
		k := AuthKey(auth.ClientId, auth.Username)
		v, err := AuthValue(auth)
		if err != nil {
			return err
		}
		return txn.Set(k, v)
	})
}

func (ar AuthBadgerRepository) DeleteByClientIdUsername(clientId string, username string) error {
	return ar.Db.Update(func(txn *badger.Txn) error {
		k := AuthKey(clientId, username)
		return txn.Delete(k)
	})
}

func (ar AuthBadgerRepository) GetAll() []model.Auth {
	auths := []model.Auth{}
	ar.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				a, err := model.GobDecode[model.Auth](v)
				if err != nil {
					return err
				}
				auths = append(auths, a)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return auths
}
