package mqtt

import (
	"strings"
)

type Subscriptions interface {
	addSub(string, string) error
	remSub(string, string) int
	findSubs(string) []string
	findSub(string) []string
}

type inMemorySubscriptions map[string][]string

func (is inMemorySubscriptions) addSub(subscribed string, subscriber string) error {
	if subs, ok := is[subscribed]; ok {
		if i := findSubscriber(subs, subscriber); i == -1 {
			is[subscribed] = append(is[subscribed], subscriber)
		}
	} else {
		is[subscribed] = []string{subscriber}
	}
	return nil
}

func (is inMemorySubscriptions) remSub(subscribed string, subscriber string) int {
	subs := is.findSub(subscribed)
	if subs != nil {
		toRem := findSubscriber(subs, subscriber)
		if toRem != -1 {
			is[subscribed] = append(is[subscribed][:toRem], is[subscribed][toRem+1:]...)
			return toRem
		}
	}
	return -1
}

func findSubscriber(subscribers []string, subscriber string) int {
	if len(subscribers) > 0 {
		for i, c := range subscribers {
			if c == subscriber {
				return i
			}
		}
	}
	return -1
}

func (is inMemorySubscriptions) findSubs(subscribed string) []string {
	// only topics have topic separator in name
	topicSegments := strings.Split(subscribed, TOPIC_SEPARATOR)
	if len(topicSegments) == 1 {
		return is.findSub(subscribed)
	} else {
		return is.multiSegmentSubs(topicSegments)
	}
}

func (is inMemorySubscriptions) findSub(subscribed string) []string {
	if s, ok := is[subscribed]; ok {
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
