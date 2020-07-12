package mqtt

import (
	"encoding/json"
	"log"

	"github.com/ilgianlu/tagyou/topic"
)

type Acl struct {
	Pattern string `json:"pattern"`
}

func CheckAcl(t string, aclsAsString string) bool {
	acls, err := readAcls(aclsAsString)
	if err != nil {
		log.Println("error checking subscribe acls :", err)
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
