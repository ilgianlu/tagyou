package mqtt

import "time"

type Subscription struct {
	clientId             string
	topic                string
	subRetainHandling    uint8
	subRetainAsPublished uint8
	subNoLocal           uint8
	subQoS               uint8
	enabled              bool
	createdAt            time.Time
}
