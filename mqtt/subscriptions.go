package mqtt

import (
	"fmt"
	"strings"
)

type Subscriptions interface {
	addSub(string, string) error
	remSub(string, string) error
	findSubs(string) []string
	findSub(string) []string
}

type inMemorySubscriptions map[string][]string

func (is inMemorySubscriptions) addSub(topic string, clientId string) error {
	if subs, ok := is[topic]; ok {
		if i := findClientid(subs, clientId); i == -1 {
			is[topic] = append(is[topic], clientId)
		}
	} else {
		is[topic] = []string{clientId}
	}
	return nil
}

func (is inMemorySubscriptions) remSub(topic string, clientId string) error {
	subs := is.findSubs(topic)
	toRem := findClientid(subs, clientId)
	if toRem != -1 {
		is[topic] = append(is[topic][:toRem], is[topic][toRem+1:]...)
		return nil
	}
	return fmt.Errorf("could not find %s in %s\n", clientId, topic)
}

func findClientid(subs []string, clientId string) int {
	if len(subs) > 0 {
		for i, c := range subs {
			if c == clientId {
				return i
			}
		}
	}
	return -1
}

func (is inMemorySubscriptions) findSubs(topic string) []string {
	topicSegments := strings.Split(topic, TOPIC_SEPARATOR)
	if len(topicSegments) == 1 {
		return is.findSub(topic)
	} else {
		return is.multiSegmentSubs(topicSegments)
	}
}

func (is inMemorySubscriptions) findSub(topic string) []string {
	if s, ok := is[topic]; ok {
		return s
	}
	return nil
}

func (is inMemorySubscriptions) multiSegmentSubs(topicSegments []string) []string {
	subs := make([]string, 0)
	for i := 1; i <= len(topicSegments); i++ {
		subT := append(make([]string, 0), topicSegments[:i]...)
		if len(subT) < len(topicSegments) {
			subT = append(subT, TOPIC_WILDCARD)
		}
		t := strings.Join(subT, TOPIC_SEPARATOR)
		ss := is.findSub(t)
		subs = append(subs, ss...)
	}
	return subs
}
