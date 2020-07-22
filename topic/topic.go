package topic

import "strings"

const TOPIC_SEPARATOR = "/"
const TOPIC_WILDCARD = "#"

func Match(topic string, matcher string) bool {
	if matcher == "#" || topic == matcher {
		return true
	}
	topicRoad := strings.Split(topic, TOPIC_SEPARATOR)
	matcherRoad := strings.Split(matcher, TOPIC_SEPARATOR)
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

func MatcherSubset(subSet string, set string) bool {
	if set == "#" {
		return true
	}
	if set == "" {
		return false
	}
	subsetRoad := strings.Split(subSet, TOPIC_SEPARATOR)
	setRoad := strings.Split(set, TOPIC_SEPARATOR)
	if len(setRoad) > len(subsetRoad) {
		return false
	}
	if len(setRoad) < len(subsetRoad) && setRoad[len(setRoad)-1] != "#" {
		return false
	}
	for i := 0; i < len(setRoad); i++ {
		if setRoad[i] != "+" && setRoad[i] != "#" && setRoad[i] != subsetRoad[i] {
			return false
		}
	}
	return true
}
