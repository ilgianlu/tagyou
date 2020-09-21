package topic

import (
	"strings"
)

const TOPIC_SEPARATOR = "/"
const TOPIC_WILDCARD = "#"
const TOPIC_JOLLY = "+"

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

func Explode(topic string) []string {
	setRoad := strings.Split(topic, TOPIC_SEPARATOR)
	res := []string{
		"#",
		topic,
	}
	if len(setRoad) == 1 {
		res = append(res, TOPIC_JOLLY)
		return res
	}
	for i := 0; i < len(setRoad); i++ {
		prev := pre(setRoad, i)
		pos := post(setRoad, i)
		if prev == "" {
			res = append(res, setRoad[i]+TOPIC_SEPARATOR+TOPIC_WILDCARD)
			res = append(res, TOPIC_JOLLY+TOPIC_SEPARATOR+TOPIC_WILDCARD)
			res = append(res, TOPIC_JOLLY+TOPIC_SEPARATOR+pos)
		} else {
			if pos == "" {
				res = append(res, prev+TOPIC_SEPARATOR+TOPIC_JOLLY)
			} else {
				res = append(res, prev+TOPIC_SEPARATOR+setRoad[i]+TOPIC_SEPARATOR+TOPIC_WILDCARD)
				res = append(res, prev+TOPIC_SEPARATOR+TOPIC_JOLLY+TOPIC_SEPARATOR+pos)
				res = append(res, prev+TOPIC_SEPARATOR+TOPIC_JOLLY+TOPIC_SEPARATOR+TOPIC_WILDCARD)
			}
		}
	}
	return res
}

func pre(path []string, i int) string {
	if i == 0 {
		return ""
	}
	return strings.Join(path[0:i], TOPIC_SEPARATOR)
}

func post(path []string, i int) string {
	if i == len(path)-1 {
		return ""
	}
	return strings.Join(path[i+1:], TOPIC_SEPARATOR)
}
