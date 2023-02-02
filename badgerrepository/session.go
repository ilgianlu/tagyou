package badgerrepository

import (
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/model"
)

type SessionBadgerRepository struct {
	Db *badger.DB
}

func SessionKey(clientId string) []byte {
	return []byte(clientId)
}

func SessionValue(session model.Session) ([]byte, error) {
	return model.GobEncode(session)
}

func (sr SessionBadgerRepository) sessionSave(session model.Session) error {
	return sr.Db.Update(func(txn *badger.Txn) error {
		k := SessionKey(session.ClientId)
		v, err := SessionValue(session)
		if err != nil {
			return err
		}
		return txn.Set(k, v)
	})
}

func (sr SessionBadgerRepository) PersistSession(running *model.RunningSession, connected bool) (sessionId uint, err error) {
	running.Mu.RLock()
	defer running.Mu.RUnlock()
	sess := model.Session{
		LastSeen:        running.LastSeen,
		LastConnect:     running.LastConnect,
		ExpiryInterval:  running.ExpiryInterval,
		ClientId:        running.ClientId,
		Connected:       connected,
		ProtocolVersion: running.ProtocolVersion,
	}
	saveErr := sr.sessionSave(sess)
	return sess.ID, saveErr
}

func (sr SessionBadgerRepository) CleanSession(clientId string) error {
	return sr.Db.Update(func(txn *badger.Txn) error {
		k := SessionKey(clientId)
		return txn.Delete(k)
	})
}

func (sr SessionBadgerRepository) SessionExists(clientId string) (model.Session, bool) {
	var session model.Session
	key := SessionKey(clientId)
	err := sr.Db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		item.Value(func(val []byte) error {
			sess, err := model.GobDecode[model.Session](val)
			if err != nil {
				return err
			}
			session = sess
			return nil
		})
		return nil
	})

	return session, err == nil
}

func (sr SessionBadgerRepository) DisconnectSession(clientId string) {
	sr.Db.View(func(txn *badger.Txn) error {
		key := SessionKey(clientId)
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		var session model.Session
		item.Value(func(val []byte) error {
			session, err := model.GobDecode[model.Session](val)
			if err != nil {
				return err
			}
			session.Connected = false
			session.LastSeen = time.Now().Unix()
			return nil
		})
		value, err := SessionValue(session)
		err = txn.Set(key, value)
		return nil
	})
}

func (sr SessionBadgerRepository) GetAll() []model.Session {
	sessions := []model.Session{}
	sr.Db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			item.Value(func(val []byte) error {
				sess, err := model.GobDecode[model.Session](val)
				if err != nil {
					return err
				}
				sessions = append(sessions, sess)
				return nil
			})
		}
		return nil
	})
	return sessions
}

func (sr SessionBadgerRepository) Save(session *model.Session) {
	sr.sessionSave(*session)
}

func (sr SessionBadgerRepository) IsOnline(clientId string) bool {
	var res bool
	sr.Db.View(func(txn *badger.Txn) error {
		key := SessionKey(clientId)
		item, err := txn.Get(key)
		if err != nil {
			return err
		}
		var session model.Session
		item.Value(func(val []byte) error {
			session, err := model.GobDecode[model.Session](val)
			if err != nil {
				return err
			}
			res = session.Connected
			return nil
		})
		value, err := SessionValue(session)
		err = txn.Set(key, value)
		return nil
	})
	return res
}
