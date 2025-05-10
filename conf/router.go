package conf

// Router modes:
// "standard" router acts as standard mqtt router
const ROUTER_MODE_STANDARD = "standard"

// "simple" stripped down router, match only fulltext topics, ignores pounds and pluses special chars
const ROUTER_MODE_SIMPLE = "simple"

// "debug" act as "standard" but also dump all published messages to debug file,
// use DEBUG_CLIENTS to dump only some clients traffic
const ROUTER_MODE_DEBUG = "debug"

// list of selected clientId to debug traffic separated by slash "/" "client1/client2", empty string dump all traffic
// if at least one client is selected then traffic dump will be separated in one file per device "client1.dump"
var DEBUG_CLIENTS = ""

// current router mode (can be overridden by env var)
var ROUTER_MODE string = ROUTER_MODE_STANDARD

// with router mode "debug": use this file to write traffic
var DEBUG_DATA_PATH string = "./debug"

// initial capacity of router (affects memory and performance at startup)
// affects only simple router and standard router
var ROUTER_STARTING_CAPACITY = 10000
