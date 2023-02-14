package conf

import (
	"os"
	"strconv"
)

func Loader() {
	FORBID_ANONYMOUS_LOGIN = os.Getenv("FORBID_ANONYMOUS_LOGIN") == "true"
	ACL_ON = os.Getenv("ACL_ON") == "true"
	CLEAN_EXPIRED_SESSIONS = os.Getenv("CLEAN_EXPIRED_SESSIONS") == "true"
	CLEAN_EXPIRED_RETRIES = os.Getenv("CLEAN_EXPIRED_RETRIES") == "true"

	var s string
	s = os.Getenv("DB_PATH")
	if s != "" {
		DB_PATH = s
	}

	s = os.Getenv("DB_NAME")
	if s != "" {
		DB_NAME = s
	}

	s = os.Getenv("CLEAN_EXPIRED_SESSIONS_INTERVAL")
	if s != "" {
		ces, err := strconv.Atoi(s)
		if err != nil {
			CLEAN_EXPIRED_SESSIONS_INTERVAL = ces
		}
	}

	s = os.Getenv("CLEAN_EXPIRED_RETRIES_INTERVAL")
	if s != "" {
		cer, err := strconv.Atoi(s)
		if err != nil {
			CLEAN_EXPIRED_RETRIES_INTERVAL = cer
		}
	}

	s = os.Getenv("DEFAULT_KEEPALIVE")
	if s != "" {
		dk, err := strconv.Atoi(s)
		if err != nil {
			DEFAULT_KEEPALIVE = dk
		}
	}

	s = os.Getenv("RETRY_EXPIRATION")
	if s != "" {
		dk, err := strconv.Atoi(s)
		if err != nil {
			RETRY_EXPIRATION = dk
		}
	}
}
