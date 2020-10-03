package conf

// (secs) session can't last more than
var SESSION_MAX_DURATION_SECONDS uint32 = 3600

// clean expired sessions async (client_id never that never reconnects)
var CLEAN_EXPIRED_SESSIONS bool = true

// (minutes) interval at which check for expired sessions cleanup (can be heavy on db, don't be too aggressive)
var CLEAN_EXPIRED_SESSIONS_INTERVAL int = 60

// (secs) keepalive for tcp connections before mqtt connect
var DEFAULT_KEEPALIVE int = 30

const LOCALHOST = "127.0.0.1"
