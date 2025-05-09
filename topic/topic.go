package topic

import (
	"strings"

	"github.com/ilgianlu/tagyou/conf"
)

// Match does topic matches the matcher (ie subscription?)
func Match(topic string, matcher string) bool {
	if matcher == "#" || topic == matcher {
		return true
	}
	topicRoad := strings.Split(topic, conf.LEVEL_SEPARATOR)
	matcherRoad := strings.Split(matcher, conf.LEVEL_SEPARATOR)
	if len(matcherRoad) > len(topicRoad) {
		return false
	}
	if len(matcherRoad) < len(topicRoad) && matcherRoad[len(matcherRoad)-1] != "#" {
		return false
	}
	for i := 0; i < len(matcherRoad); i++ {
		if matcherRoad[i] != "+" && matcherRoad[i] != "#" && matcherRoad[i] != topicRoad[i] {
			return false
		}
	}
	return true
}

// MatcherSubset is subMatcher included in matcher (acl, can client subscribe this?)
func MatcherSubset(subMatcher string, matcher string) bool {
	if matcher == "#" {
		return true
	}
	if matcher == "" {
		return false
	}
	subMatcherRoad := strings.Split(subMatcher, conf.LEVEL_SEPARATOR)
	setRoad := strings.Split(matcher, conf.LEVEL_SEPARATOR)
	if len(setRoad) > len(subMatcherRoad) {
		return false
	}
	if len(setRoad) < len(subMatcherRoad) && setRoad[len(setRoad)-1] != "#" {
		return false
	}
	for i, c := range setRoad {
		if c != "+" && c != "#" && c != subMatcherRoad[i] {
			return false
		}
	}
	return true
}

func SharedSubscription(topic string) bool {
	if len(topic) <= 10 {
		// min shared topic length
		return false
	}
	return topic[0:6] == conf.TOPIC_SHARED && string(topic[7]) != conf.LEVEL_SEPARATOR
}

// SharedSubscriptionTopicParse return if shared
func SharedSubscriptionTopicParse(topic string) (string, string) {
	s := topic[len(conf.TOPIC_SHARED)+1:]
	i := strings.Index(s, conf.LEVEL_SEPARATOR)
	shareName := s[:i]
	subscribedTopic := s[i+1:]
	return shareName, subscribedTopic
}
