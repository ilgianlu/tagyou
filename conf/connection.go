package conf

// (secs) session can't last more than
var SESSION_MAX_DURATION_SECONDS uint32 = 3600

// clean expired sessions async (client_id never that never reconnects)
var CLEAN_EXPIRED_SESSIONS bool = true

// (minutes) interval at which check for expired sessions cleanup (can be heavy on db, don't be too aggressive)
var CLEAN_EXPIRED_SESSIONS_INTERVAL int = 60

// (secs) keepalive for tcp connections before mqtt connect
var DEFAULT_KEEPALIVE int = 30

// (qos > 0) stop retries (wait for) for messages pub_ack pub_rec pub_rel pub_comp
var CLEAN_EXPIRED_RETRIES bool = true

// (minutes) interval at which check for expired retries cleanup (can be heavy on db, don't be too aggressive)
var CLEAN_EXPIRED_RETRIES_INTERVAL int = 60

// (secs) consider retry expired after RETRY_EXPIRATION secs
var RETRY_EXPIRATION int = 60
