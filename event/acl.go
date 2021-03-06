package event

import (
	"encoding/json"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/topic"
)

type Acl struct {
	Pattern string `json:"pattern"`
}

func CheckAcl(t string, aclsAsString string) bool {
	if aclsAsString == "" {
		return false
	}
	acls, err := readAcls(aclsAsString)
	if err != nil {
		log.Err(err).Msg("error checking subscribe acls")
		return false
	}
	for _, acl := range acls {
		if topic.MatcherSubset(t, acl.Pattern) {
			return true
		}
	}
	return false
}

func readAcls(aclsAsString string) ([]Acl, error) {
	var acls []Acl
	if err := json.Unmarshal([]byte(aclsAsString), &acls); err != nil {
		return acls, err
	}
	return acls, nil
}
