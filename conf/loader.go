package conf

import (
	"os"
	"strconv"
)

func Loader() {
	FORBID_ANONYMOUS_LOGIN = os.Getenv("FORBID_ANONYMOUS_LOGIN") != "false"
	ACL_ON = os.Getenv("ACL_ON") == "true"
	CLEAN_EXPIRED_SESSIONS = os.Getenv("CLEAN_EXPIRED_SESSIONS") != "false"
	CLEAN_EXPIRED_RETRIES = os.Getenv("CLEAN_EXPIRED_RETRIES") != "false"
	INIT_DB = os.Getenv("INIT_DB") == "true" || os.Getenv("INIT_DB") == "1"

	var s string
	s = os.Getenv("API_PORT")
	if s != "" {
		API_PORT = s
	}

	s = os.Getenv("WS_PORT")
	if s != "" {
		WS_PORT = s
	}

	s = os.Getenv("LISTEN_PORT")
	if s != "" {
		LISTEN_PORT = s
	}

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

	s = os.Getenv("API_TOKEN_SIGNING_KEY")
	if s != "" {
		API_TOKEN_SIGNING_KEY = []byte(s)
	}

	s = os.Getenv("API_TOKEN_ISSUER")
	if s != "" {
		API_TOKEN_ISSUER = s
	}

	s = os.Getenv("API_TOKEN_HOURS_DURATION")
	if s != "" {
		dk, err := strconv.Atoi(s)
		if err != nil {
			API_TOKEN_HOURS_DURATION = dk
		}
	}

	s = os.Getenv("INIT_ADMIN_PASSWORD")
	if s != "" {
		INIT_ADMIN_PASSWORD = []byte(s)
	}

	s = os.Getenv("ROUTER_MODE")
	if s != "" {
		ROUTER_MODE = s
	}

	s = os.Getenv("DEBUG_CLIENTS")
	if s != "" {
		DEBUG_CLIENTS = s
	}

	s = os.Getenv("SIMPLE_CLIENTS")
	if s != "" {
		SIMPLE_CLIENTS = s
	}

	s = os.Getenv("DEBUG_DATA_PATH")
	if s != "" {
		DEBUG_DATA_PATH = s
	}

	s = os.Getenv("AI_URL")
	if s != "" {
		AI_URL = s
	}

	s = os.Getenv("AI_MODEL")
	if s != "" {
		AI_MODEL = s
	}

	s = os.Getenv("AI_API_KEY")
	if s != "" {
		AI_API_KEY = s
	}
}
