package nowherecloud

func respectFilter(topic string) (respectedFilter string, found bool) {
	for _, t := range KAFKA_TOPICS {
		if len(topic) <= len(t) {
			continue
		}
		if topic[:len(t)] == t {
			return t, true
		}
	}
	return "", false
}
