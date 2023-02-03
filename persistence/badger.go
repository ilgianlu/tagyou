package persistence

import (
	"encoding/gob"

	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/badgerrepository"
	"github.com/ilgianlu/tagyou/model"
)

type BadgerPersistence struct {
	dba *badger.DB
	dbb *badger.DB
	dbc *badger.DB
	dbd *badger.DB
	dbe *badger.DB
}

var (
	badgerPersistence BadgerPersistence
)

func (p BadgerPersistence) Init() error {
	gob.Register(model.Auth{})
	gob.Register(model.Retain{})
	gob.Register(model.Retry{})
	gob.Register(model.Session{})
	gob.Register(model.Subscription{})

	dba, _ := badger.Open(badger.DefaultOptions("db/auth.db"))
	AuthRepository = badgerrepository.AuthBadgerRepository{Db: dba}
	badgerrepository.StartGarbageCollection(dba)

	dbb, _ := badger.Open(badger.DefaultOptions("db/retain.db"))
	RetainRepository = badgerrepository.RetainBadgerRepository{Db: dbb}
	badgerrepository.StartGarbageCollection(dbb)

	dbc, _ := badger.Open(badger.DefaultOptions("db/retry.db"))
	RetryRepository = badgerrepository.RetryBadgerRepository{Db: dbc}
	badgerrepository.StartGarbageCollection(dbc)

	dbd, _ := badger.Open(badger.DefaultOptions("db/session.db"))
	SessionRepository = badgerrepository.SessionBadgerRepository{Db: dbd}
	badgerrepository.StartSessionCleaner(dbd)
	badgerrepository.StartGarbageCollection(dbd)

	dbe, _ := badger.Open(badger.DefaultOptions("db/subscription.db"))
	SubscriptionRepository = badgerrepository.SubscriptionBadgerRepository{Db: dbe}
	badgerrepository.StartGarbageCollection(dbe)

	badgerPersistence = BadgerPersistence{
		dba: dba,
		dbb: dbb,
		dbc: dbc,
		dbd: dbd,
		dbe: dbe,
	}

	return nil
}

func (p BadgerPersistence) Close() {
	p.dba.Close()
	p.dbb.Close()
	p.dbc.Close()
	p.dbd.Close()
	p.dbe.Close()
}
