package conf

// Router modes:
// "standard" router acts as standard mqtt router
// "simple" stripped down router, match only fulltext topics, ignores pounds and pluses special chars
// "debug" act as "standard" but also dump all published messages to debug file
var ROUTER_MODE string = "standard"

// with router mode "debug": use this file to write traffic
var DEBUG_FILE string = "traffic.debug"
