package conf

// Router modes:
// "standard" router acts as standard mqtt router
const ROUTER_MODE_STANDARD = "standard"

// "simple" stripped down router, match only fulltext topics, ignores pounds and pluses special chars
const ROUTER_MODE_SIMPLE = "simple"

// "debug" act as "standard" but also dump all published messages to debug file (one per client id)
const ROUTER_MODE_DEBUG = "debug"

// list of selected clientId to use for debug router separated by slash "/" "client1/client2", empty string use default 
var DEBUG_CLIENTS = ""

// list of selected clientId to use for simple router separated by slash "/" "client1/client2", empty string use default 
var SIMPLE_CLIENTS = ""

// current router mode (can be overridden by env var)
var ROUTER_MODE string = ROUTER_MODE_STANDARD

// with router mode "debug": use this file to write traffic
var DEBUG_DATA_PATH string = "./debug"

// initial capacity of router (affects memory and performance at startup)
// affects only simple router and standard router
var ROUTER_STARTING_CAPACITY = 10000
