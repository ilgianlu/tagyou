package badgerrepository

import (
	"fmt"

	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/robfig/cron"
)

func StartGarbageCollection(db *badger.DB) {
	log.Info().Msg("[MQTT] start badger garbage collection")
	c := cron.New()
	_ = c.AddFunc(
		fmt.Sprintf("@every %dm", conf.CLEAN_EXPIRED_RETRIES_INTERVAL),
		func() {
			log.Info().Msg("[MQTT] collecting garbage")
			db.RunValueLogGC(0.7)
		},
	)
	c.Start()
}
