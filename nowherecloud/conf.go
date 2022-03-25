package nowherecloud

import (
	"os"
	"strings"
)

var KAFKA_ON bool = false
var KAFKA_URL string = "localhost:9092"
var KAFKA_TOPICS []string = []string{"topic_one", "topic_two"}

var REDIS_URL string = "localhost:6379"

var BROKER_PROTOCOL string = "mqtt"
var BROKER_PORT string = "1883"

func Loader() {
	KAFKA_ON = os.Getenv("KAFKA_ON") == "true"
	KAFKA_URL = os.Getenv("KAFKA_URL")
	KAFKA_TOPICS = strings.Split(os.Getenv("KAFKA_TOPICS"), ",")
	if os.Getenv("REDIS_URL") != "" {
		REDIS_URL = os.Getenv("REDIS_URL")
	}
	if os.Getenv("BROKER_PROTOCOL") != "" {
		BROKER_PROTOCOL = os.Getenv("BROKER_PROTOCOL")
	}
	if os.Getenv("BROKER_PORT") != "" {
		BROKER_PROTOCOL = os.Getenv("BROKER_PORT")
	}
}
