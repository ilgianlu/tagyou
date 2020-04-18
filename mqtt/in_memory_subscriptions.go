package mqtt

import (
	"fmt"
	"log"
	"strings"
)

type inMemorySubscriptions struct {
	clientSubscriptions map[string][]string
	topicSubscriptions  map[string][]string
}

func (is inMemorySubscriptions) addSubscription(topic string, subscriber string) error {
	addIn(&is.topicSubscriptions, topic, subscriber)
	addIn(&is.clientSubscriptions, subscriber, topic)
	return nil
}

func (is inMemorySubscriptions) remSubscription(topic string, subscriber string) error {
	res1 := removeIn(&is.topicSubscriptions, topic, subscriber)
	res0 := removeIn(&is.clientSubscriptions, subscriber, topic)
	if res0 != -1 && res1 != -1 {
		return nil
	} else {
		return fmt.Errorf("could not remove subscription (%d, %d)", res0, res1)
	}
}

func (is inMemorySubscriptions) findSubscribers(topic string) []string {
	// only topics have topic separator in name
	topicSegments := strings.Split(topic, TOPIC_SEPARATOR)
	if len(topicSegments) == 1 {
		if s, ok := is.findSubscribed(topic); ok {
			return s
		} else {
			return []string{}
		}
	} else {
		return is.multiSegmentSubs(topicSegments)
	}
}

func (is inMemorySubscriptions) findSubscribed(topic string) ([]string, bool) {
	s, ok := is.topicSubscriptions[topic]
	return s, ok
}

func (is inMemorySubscriptions) remSubscribed(topic string) {
	delete(is.topicSubscriptions, topic)
}

func (is inMemorySubscriptions) multiSegmentSubs(topicSegments []string) []string {
	subs := make([]string, 0)
	for i := 1; i <= len(topicSegments); i++ {
		subT := append(make([]string, 0), topicSegments[:i]...)
		if len(subT) < len(topicSegments) {
			subT = append(subT, TOPIC_WILDCARD)
		}
		t := strings.Join(subT, TOPIC_SEPARATOR)
		if ss, ok := is.findSubscribed(t); ok {
			subs = append(subs, ss...)
		}
	}
	return subs
}

func findIn(subscribers []string, subscriber string) int {
	if len(subscribers) > 0 {
		for i, c := range subscribers {
			if c == subscriber {
				return i
			}
		}
	}
	return -1
}

func addIn(sub *map[string][]string, key string, val string) {
	if subs, ok := (*sub)[key]; ok {
		if i := findIn(subs, val); i == -1 {
			(*sub)[key] = append((*sub)[key], val)
		}
	} else {
		(*sub)[key] = []string{val}
	}
}

func removeIn(sub *map[string][]string, key string, val string) int {
	if vals, ok := (*sub)[key]; ok {
		log.Println("removing", val, "from", key)
		toRem := findIn(vals, val)
		if toRem != -1 {
			(*sub)[key] = append((*sub)[key][:toRem], (*sub)[key][toRem+1:]...)
		}
		return toRem
	} else {
		log.Println("could not remove", val, "from", vals, "not present")
		return -1
	}
}
