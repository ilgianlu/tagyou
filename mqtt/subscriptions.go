package mqtt

type Subscriptions interface {
	addSubscription(string, string) error
	remSubscription(string, string) error
	findSubscribers(string) []string
	findSubscribed(string) ([]string, bool)
	remSubscribed(string)
}
