package sqlrepository

import (
	"fmt"
	"log/slog"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func StartSessionCleaner(db *gorm.DB) {
	slog.Info("[MQTT] start expired sessions cleaner")
	cleanSessions(db)
	c := cron.New()
	_ = c.AddFunc(
		fmt.Sprintf("@every %dm", conf.CLEAN_EXPIRED_SESSIONS_INTERVAL),
		func() {
			slog.Info("[MQTT] clean expired sessions")
			cleanSessions(db)
		},
	)
	c.Start()
}

func cleanSessions(db *gorm.DB) {
	disconnectedSessions := []Session{}
	if err := db.Debug().Where("connected = 0").Find(&disconnectedSessions); err != nil {
		for _, disconnectSession := range disconnectedSessions {
			if model.SessionExpired(disconnectSession.LastSeen, disconnectSession.ExpiryInterval) {
				db.Debug().Delete(&disconnectSession)
			}
		}
	}
	slog.Info("[MQTT] expired sessions cleanup done")
}
