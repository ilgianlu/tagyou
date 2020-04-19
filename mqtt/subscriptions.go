package mqtt

type Subscriptions interface {
	addSubscription(s Subscription) error
	remSubscription(topic string, clientId string) error
	findSubscriptionsByTopic(topic string) []Subscription
	findSubscriptionsByClientId(clientId string) []Subscription
	findTopicSubscribers(topic string) []Subscription
	remSubscriptionsByTopic(topic string)
	remSubscriptionsByClientId(clientId string)
	disableClientSubscriptions(clientId string)
	enableClientSubscriptions(clientId string)
}
