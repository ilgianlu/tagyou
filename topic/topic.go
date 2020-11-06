package topic

import (
	"math"
	"strings"

	"github.com/ilgianlu/tagyou/conf"
)

// TopicSeparator separator between topic parts
const TopicSeparator = "/"

// TopicWildcard matches multiple consecutive parts in the end of a topic
const TopicWildcard = "#"

// TopicJolly matches a single part in a topic "road"
const TopicJolly = "+"

// Match does topic matches the matcher (ie subscription?)
func Match(topic string, matcher string) bool {
	if matcher == "#" || topic == matcher {
		return true
	}
	topicRoad := strings.Split(topic, TopicSeparator)
	matcherRoad := strings.Split(matcher, TopicSeparator)
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
	subMatcherRoad := strings.Split(subMatcher, TopicSeparator)
	setRoad := strings.Split(matcher, TopicSeparator)
	if len(setRoad) > len(subMatcherRoad) {
		return false
	}
	if len(setRoad) < len(subMatcherRoad) && setRoad[len(setRoad)-1] != "#" {
		return false
	}
	for i := 0; i < len(setRoad); i++ {
		if setRoad[i] != "+" && setRoad[i] != "#" && setRoad[i] != subMatcherRoad[i] {
			return false
		}
	}
	return true
}

// Explode list all possible subscriptions to look for when client publish a message in a topic
func Explode(topic string) []string {
	if conf.Matcher == conf.MatcherBasic {
		return []string{topic}
	}
	road := strings.Split(topic, TopicSeparator)
	if conf.Matcher == conf.MatcherMultilevelOnly {
		return explodeMultiLevel(road)
	}
	return explodeFull(road)
}

func explodeMultiLevel(road []string) []string {
	res := []string{}
	for i := 0; i < len(road); i++ {
		r := append([]string{}, road[:i]...)
		r = append(r, TopicWildcard)
		t := strings.Join(r, TopicSeparator)
		res = append(res, t)
	}
	return res
}

func explodeFull(road []string) []string {
	res := []string{"#"}
	for i := 1; i <= len(road); i++ {
		subRoads := explodeSingleLevel(road[:i])
		for _, subRoad := range subRoads {
			if i != len(road) {
				subRoad = append(subRoad, TopicWildcard)
			}
			res = append(res, strings.Join(subRoad, "/"))
		}
	}
	return res
}

func explodeSingleLevel(road []string) [][]string {
	res := [][]string{}
	l := math.Pow(2, float64(len(road)))
	for i := 0; i < int(l); i++ {
		res = append(res, singleLevel(road, i))
	}
	return res
}

func singleLevel(road []string, i int) []string {
	ss := []string{}
	for p, e := range road {
		o := i & (1 << p)
		if o > 0 {
			ss = append(ss, TopicJolly)
		} else {
			ss = append(ss, e)
		}
	}
	return ss
}
