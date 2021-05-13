package conf

import (
	"os"
	"strconv"
)

func Loader() {
	FORBID_ANONYMOUS_LOGIN = os.Getenv("FORBID_ANONYMOUS_LOGIN") == "true"
	ACL_ON = os.Getenv("ACL_ON") == "true"
	CLEAN_EXPIRED_SESSIONS = os.Getenv("CLEAN_EXPIRED_SESSIONS") == "true"
	KAFKA_ON = os.Getenv("KAFKA_ON") == "true"

	var s string
	s = os.Getenv("CLEAN_EXPIRED_SESSIONS_INTERVAL")
	if s != "" {
		ces, err := strconv.Atoi(s)
		if err != nil {
			CLEAN_EXPIRED_SESSIONS_INTERVAL = ces
		}
	}

	s = os.Getenv("DEFAULT_KEEPALIVE")
	if s != "" {
		dk, err := strconv.Atoi(s)
		if err != nil {
			DEFAULT_KEEPALIVE = dk
		}
	}
}
