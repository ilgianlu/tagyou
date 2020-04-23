package mqtt

type Retains interface {
	addRetain(retain Retain) error
	remRetain(topic string) error
	findRetainsByTopic(topic string) []Retain
}
