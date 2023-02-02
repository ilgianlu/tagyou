package persistence

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/ilgianlu/tagyou/badgerrepository"
)

func InitBadgerRepositories() {
	dba, _ := badger.Open(badger.DefaultOptions("db/auth.db"))
	AuthRepository = badgerrepository.AuthBadgerRepository{Db: dba}
	dbb, _ := badger.Open(badger.DefaultOptions("db/retain.db"))
	RetainRepository = badgerrepository.RetainBadgerRepository{Db: dbb}
	dbc, _ := badger.Open(badger.DefaultOptions("db/retry.db"))
	RetryRepository = badgerrepository.RetryBadgerRepository{Db: dbc}
	dbd, _ := badger.Open(badger.DefaultOptions("db/session.db"))
	SessionRepository = badgerrepository.SessionBadgerRepository{Db: dbd}
	dbe, _ := badger.Open(badger.DefaultOptions("db/subscription.db"))
	SubscriptionRepository = badgerrepository.SubscriptionBadgerRepository{Db: dbe}
}
