package mqtt

type Retains interface {
	addRetain(retain Retain) error
	remRetain(topic string) error
	findRetainByTopic(topic string) []Retain
}
