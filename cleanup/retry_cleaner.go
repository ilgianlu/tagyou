package cleanup

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/sqlrepository"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func StartRetryCleaner(db *gorm.DB) {
	log.Info().Msg("[MQTT] start expired retries cleaner")
	cleanRetries(db)
	c := cron.New()
	_ = c.AddFunc(
		fmt.Sprintf("@every %dm", conf.CLEAN_EXPIRED_RETRIES_INTERVAL),
		func() {
			log.Info().Msg("[MQTT] clean expired retries")
			cleanRetries(db)
		},
	)
	c.Start()
}

func cleanRetries(db *gorm.DB) {
	expireTime := time.Now().Unix() - int64(conf.RETRY_EXPIRATION)
	if err := db.Debug().Where("created_at < ?", expireTime).Delete(&sqlrepository.Retry{}).Error; err != nil {
		log.Error().Err(err)
	}
	log.Info().Msg("[MQTT] retries cleanup done")
}
