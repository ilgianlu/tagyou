package model

type Packet interface {
	QoS() byte
	Retain() bool
	ApplicationMessage() []byte
	PacketIdentifier() int
	Dup() bool
	GetPublishTopic() string
	GetReasonCode() uint8
	GetSubscriptions() []Subscription
	PacketType() byte
}
