package sqlrepository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
	"github.com/robfig/cron"
)

func StartRetryCleaner(db *dbaccess.Queries) {
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

func cleanRetries(db *dbaccess.Queries) {
	expireTime := time.Now().Unix() - int64(conf.RETRY_EXPIRATION)
	err := db.DeleteRetriesOlder(context.Background(), sql.NullInt64{Int64: expireTime, Valid: true})
	if err != nil {
		slog.Error("", "err", err)
	}
	slog.Info("[MQTT] retries cleanup done")
}
