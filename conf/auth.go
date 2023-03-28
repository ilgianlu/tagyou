package conf

var FORBID_ANONYMOUS_LOGIN bool = true
var ACL_ON bool = true
var PASSWORD_MIN_LENGTH int = 8

var API_TOKEN_SIGNING_KEY []byte = []byte("generic_non_acceptable_key_2023")
var API_TOKEN_ISSUER string = "TagYou"
var API_TOKEN_HOURS_DURATION int = 1
