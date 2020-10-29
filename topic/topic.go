package topic

import (
	"math"
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

	// res = append(res, explodeSingleLevel(road)...)
	// res = append(res, explodeMultiLevel(road)...)
	if len(setRoad) == 1 {
		res = append(res, TOPIC_JOLLY)
		return res
	}
	for i := 0; i < len(setRoad); i++ {
		prev, isPrev := pre(setRoad, i)
		pos, isPos := post(setRoad, i)
		if !isPrev {
			res = append(res, setRoad[i]+TOPIC_SEPARATOR+TOPIC_WILDCARD)
			res = append(res, TOPIC_JOLLY+TOPIC_SEPARATOR+TOPIC_WILDCARD)
			res = append(res, TOPIC_JOLLY+TOPIC_SEPARATOR+pos)
		} else {
			if !isPos {
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

func explodeMultiLevel(road []string) []string {
	res := []string{}
	for i := 0; i < len(road); i++ {
		r := append([]string{}, road[:i]...)
		r = append(r, TOPIC_WILDCARD)
		t := strings.Join(r, TOPIC_SEPARATOR)
		res = append(res, t)
	}
	return res
}

func explodeSingleLevel(road []string) []string {
	res := []string{}
	l := math.Pow(2, float64(len(road)))
	for i := 0; i < int(l); i++ {
		res = append(res, singleLevel(road, i))
	}
	return res
}

func singleLevel(road []string, i int) string {
	ss := []string{}
	for p, e := range road {
		o := i & (1 << p)
		if o > 0 {
			ss = append(ss, "+")
		} else {
			ss = append(ss, e)
		}
	}
	return strings.Join(ss, "/")
}

func pre(path []string, i int) (string, bool) {
	if i == 0 {
		return "", false
	}
	return strings.Join(path[0:i], TOPIC_SEPARATOR), true
}

func post(path []string, i int) (string, bool) {
	if i == len(path)-1 {
		return "", false
	}
	return strings.Join(path[i+1:], TOPIC_SEPARATOR), true
}
