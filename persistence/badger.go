package persistence

import (
	"encoding/gob"
	"os"

	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/badgerrepository"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

var DB_AUTH string
var DB_RETAIN string
var DB_RETRY string
var DB_SESSION string
var DB_SUBSCRIPTION string

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
	DB_AUTH = conf.DB_PATH + "auth.db"
	DB_RETAIN = conf.DB_PATH + "retain.db"
	DB_RETRY = conf.DB_PATH + "retry.db"
	DB_SESSION = conf.DB_PATH + "session.db"
	DB_SUBSCRIPTION = conf.DB_PATH + "subscription.db"

	gob.Register(model.Auth{})
	gob.Register(model.Retain{})
	gob.Register(model.Retry{})
	gob.Register(model.Session{})
	gob.Register(model.Subscription{})

	dba, _ := badger.Open(badger.DefaultOptions(DB_AUTH))
	AuthRepository = badgerrepository.AuthBadgerRepository{Db: dba}
	badgerrepository.StartGarbageCollection(dba)

	dbb, _ := badger.Open(badger.DefaultOptions(DB_RETAIN))
	RetainRepository = badgerrepository.RetainBadgerRepository{Db: dbb}
	badgerrepository.StartGarbageCollection(dbb)

	dbc, _ := badger.Open(badger.DefaultOptions(DB_RETRY))
	RetryRepository = badgerrepository.RetryBadgerRepository{Db: dbc}
	badgerrepository.StartGarbageCollection(dbc)

	dbd, _ := badger.Open(badger.DefaultOptions(DB_SESSION))
	SessionRepository = badgerrepository.SessionBadgerRepository{Db: dbd}
	badgerrepository.StartSessionCleaner(dbd)
	badgerrepository.StartGarbageCollection(dbd)

	dbe, _ := badger.Open(badger.DefaultOptions(DB_SUBSCRIPTION))
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

func (p BadgerPersistence) CleanUp() {
	p.Close()
	os.Remove(DB_AUTH)
	os.Remove(DB_RETAIN)
	os.Remove(DB_RETRY)
	os.Remove(DB_SESSION)
	os.Remove(DB_SUBSCRIPTION)
}
