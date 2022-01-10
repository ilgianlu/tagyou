package nowherecloud

func respectFilter(topic string) bool {
	for _, t := range KAFKA_TOPICS {
		if len(topic) <= len(t) {
			continue
		}
		if topic[:len(t)] == t {
			return true
		}
	}
	return false
}
