package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/api"
	"github.com/ilgianlu/tagyou/conf"
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

	p := persistence.BadgerPersistence{}
	p.Init()
	defer p.Close()

	go api.StartApi(os.Getenv("API_PORT"))
	mq.StartMQTT(os.Getenv("LISTEN_PORT"))
}

func loadEnv() error {
	env := os.Getenv("TAGYOU_ENV")
	if env == "" {
		env = "default"
	}
	return dotenv.Load(".env." + env + ".local")
}
