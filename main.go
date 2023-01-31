package main

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/ilgianlu/tagyou/api"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	mq "github.com/ilgianlu/tagyou/mqtt"
	"github.com/ilgianlu/tagyou/persistence"
	dotenv "github.com/joho/godotenv"
)

func main() {
	err := loadEnv()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load env")
		return
	}
	conf.Loader()

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if os.Getenv("DEBUG") != "" {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	db, err := openDb()
	if err != nil {
		log.Fatal().Err(err).Msg("[MQTT] failed to connect database")
	}
	log.Info().Msg("[MQTT] db connected !")
	defer closeDb(db)

	model.Migrate(db)

	persistence.InitSqlRepositories(db)

	if conf.CLEAN_EXPIRED_SESSIONS {
		StartSessionCleaner(db)
	}
	if conf.CLEAN_EXPIRED_RETRIES {
		StartRetryCleaner(db)
	}

	go api.StartApi(os.Getenv("API_PORT"))
	mq.StartMQTT(os.Getenv("LISTEN_PORT"))
}

func loadEnv() error {
	env := os.Getenv("TAGYOU_ENV")
	if "" == env {
		env = "default"
	}
	return dotenv.Load(".env." + env + ".local")
}

func openDb() (*gorm.DB, error) {
	logLevel := logger.Silent
	if os.Getenv("DEBUG") != "" {
		logLevel = logger.Info
	}
	return gorm.Open(sqlite.Open(os.Getenv("DB_PATH")+os.Getenv("DB_NAME")), &gorm.Config{
		Logger: logger.New(
			&log.Logger,
			logger.Config{
				SlowThreshold: 200 * time.Millisecond,
				LogLevel:      logLevel,
				Colorful:      true,
			},
		),
	})
}

func closeDb(db *gorm.DB) {
	sql, err := db.DB()
	if err != nil {
		log.Error().Err(err).Msg("could not close DB")
		return
	}
	sql.Close()
}
