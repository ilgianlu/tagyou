package badgerrepository

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/model"
)

func AuthKey(clientId string, username string) []byte {
	k := []byte(clientId)
	k = append(k, []byte(username)...)
	return k
}

func AuthValue(auth model.Auth) ([]byte, error) {
	res := bytes.Buffer{}
	enc := gob.NewEncoder(&res)
	err := enc.Encode(auth)
	if err != nil {
		return []byte{}, err
	}
	return res.Bytes(), nil
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
			valReader := bytes.NewReader(val)
			decoder := gob.NewDecoder(valReader)
			err := decoder.Decode(&a)
			if err != nil {
				return err
			}
			return nil
		})

		return err
	})

	return a, err
}
