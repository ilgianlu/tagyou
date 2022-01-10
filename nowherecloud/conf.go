package nowherecloud

import (
	"os"
	"strings"
)

var KAFKA_ON bool = false
var KAFKA_URL string = "localhost:9092"
var KAFKA_TOPICS []string = []string{"topic_one", "topic_two"}

func Loader() {
	KAFKA_ON = os.Getenv("KAFKA_ON") == "true"
	KAFKA_URL = os.Getenv("KAFKA_URL")
	KAFKA_TOPICS = strings.Split(os.Getenv("KAFKA_TOPICS"), ",")

}
