package mqtt

import (
	"fmt"
	"strings"
)

type Subscriptions interface {
	addSubscription(string, string) error
	remSubscription(string, string) int
	findSubscribers(string) []string
	findSubscribed(string) ([]string, bool)
	remSubscribed(string)
}

type inMemorySubscriptions map[string][]string

func (is inMemorySubscriptions) addSubscription(subscribed string, subscriber string) error {
	if subs, ok := is[subscribed]; ok {
		if i := findSubscriber(subs, subscriber); i == -1 {
			is[subscribed] = append(is[subscribed], subscriber)
		}
	} else {
		is[subscribed] = []string{subscriber}
	}
	return nil
}

func (is inMemorySubscriptions) remSubscription(subscribed string, subscriber string) int {
	if subs, ok := is.findSubscribed(subscribed); ok {
		fmt.Println("removing", subscriber, "from", subs)
		toRem := findSubscriber(subs, subscriber)
		if toRem != -1 {
			is[subscribed] = append(is[subscribed][:toRem], is[subscribed][toRem+1:]...)
		}
		return toRem
	} else {
		fmt.Println("could not remove", subscriber, "from", subs, "not present")
		return -1
	}
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

func (is inMemorySubscriptions) findSubscribers(subscribed string) []string {
	// only topics have topic separator in name
	topicSegments := strings.Split(subscribed, TOPIC_SEPARATOR)
	if len(topicSegments) == 1 {
		if s, ok := is.findSubscribed(subscribed); ok {
			return s
		} else {
			return []string{}
		}
	} else {
		return is.multiSegmentSubs(topicSegments)
	}
}

func (is inMemorySubscriptions) findSubscribed(subscribed string) ([]string, bool) {
	s, ok := is[subscribed]
	return s, ok
}

func (is inMemorySubscriptions) remSubscribed(subscribed string) {
	delete(is, subscribed)
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
