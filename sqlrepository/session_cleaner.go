package sqlrepository

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
	"github.com/robfig/cron"
)

func StartSessionCleaner(db *dbaccess.Queries) {
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

func cleanSessions(db *dbaccess.Queries) {
	disconnectedSessions, err := db.GetDisconnectedSessions(context.Background())
	if err != nil {
		slog.Error("[MQTT] error loading sessions", "err", err)
		return
	}

	for _, disconnectSession := range disconnectedSessions {
		if model.SessionExpired(disconnectSession.LastSeen.Int64, disconnectSession.ExpiryInterval.Int64) {
			db.DeleteSessionByClientId(context.Background(), disconnectSession.ClientID)
		}
	}
	slog.Info("[MQTT] expired sessions cleanup done")
}
