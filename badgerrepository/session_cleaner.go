package badgerrepository

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/robfig/cron"
)

func StartSessionCleaner(db *badger.DB) {
	log.Info().Msg("[MQTT] start expired sessions cleaner")
	cleanSessions(db)
	c := cron.New()
	_ = c.AddFunc(
		fmt.Sprintf("@every %dm", conf.CLEAN_EXPIRED_SESSIONS_INTERVAL),
		func() {
			log.Info().Msg("[MQTT] clean expired sessions")
			cleanSessions(db)
		},
	)
	c.Start()
}

func cleanSessions(db *badger.DB) {
	db.Update(func(txn *badger.Txn) error {
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
				if sess.Expired() {
					if err := txn.Delete(item.Key()); err == badger.ErrTxnTooBig {
						_ = txn.Commit()
						txn = db.NewTransaction(true)
						_ = txn.Delete(item.Key())
					}
				}
				return nil
			})
		}
		return nil
	})
	log.Info().Msg("[MQTT] expired sessions cleanup done")
}
