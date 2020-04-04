package mqtt

type Subscriptions interface {
	addSub(string, string) error
	remSub(string, string) error
	findSubs(string) []string
}

type inMemorySubscriptions map[string][]string

func (is inMemorySubscriptions) addSub(topic string, clientId string) error {
	if _, ok := is[topic]; !ok {
		is[topic] = make([]string, 0)
	}
	is[topic] = append(is[topic], clientId)
	return nil
}

func (is inMemorySubscriptions) remSub(topic string, clientId string) error {
	return nil
}

func (is inMemorySubscriptions) findSubs(topic string) []string {
	if s, ok := is[topic]; ok {
		return s
	}
	return nil
}
