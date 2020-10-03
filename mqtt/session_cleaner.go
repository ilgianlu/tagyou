package mqtt

import (
	"fmt"
	"log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func StartSessionCLeaner(db *gorm.DB) {
	log.Println("[MQTT] start expired sessions cleaner")
	c := cron.New()
	_ = c.AddFunc(
		fmt.Sprintf("@every %dm", conf.CLEAN_EXPIRED_SESSIONS_INTERVAL),
		func() {
			log.Println("[MQTT] clean expired sessions")
			CleanSessions(db)
		},
	)
	c.Start()
}

func CleanSessions(db *gorm.DB) {
	disconnectedSessions := []model.Session{}
	if err := db.Debug().Where("connected = 0").Find(&disconnectedSessions); err != nil {
		for _, disconnectSession := range disconnectedSessions {
			if disconnectSession.Expired() {
				db.Debug().Delete(&disconnectSession)
			}
		}
	}
}
