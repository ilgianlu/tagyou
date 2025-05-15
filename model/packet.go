package model

type Packet interface {
	QoS() byte
	Retain() bool
	ApplicationMessage() []byte
}
