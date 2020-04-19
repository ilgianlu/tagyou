package mqtt

import "time"

type Subscription struct {
	clientId          string
	topic             string
	RetainHandling    uint8
	RetainAsPublished uint8
	NoLocal           uint8
	QoS               uint8
	enabled           bool
	createdAt         time.Time
}
