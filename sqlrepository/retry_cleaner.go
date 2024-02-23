package sqlrepository

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func StartRetryCleaner(db *gorm.DB) {
	slog.Info("[MQTT] start expired retries cleaner")
	cleanRetries(db)
	c := cron.New()
	_ = c.AddFunc(
		fmt.Sprintf("@every %dm", conf.CLEAN_EXPIRED_RETRIES_INTERVAL),
		func() {
			slog.Info("[MQTT] clean expired retries")
			cleanRetries(db)
		},
	)
	c.Start()
}

func cleanRetries(db *gorm.DB) {
	expireTime := time.Now().Unix() - int64(conf.RETRY_EXPIRATION)
	if err := db.Debug().Where("created_at < ?", expireTime).Delete(&model.Retry{}).Error; err != nil {
		slog.Error("", "err", err)
	}
	slog.Info("[MQTT] retries cleanup done")
}
