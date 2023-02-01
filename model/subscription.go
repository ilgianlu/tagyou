package model

type Subscription struct {
	ClientId          string
	Topic             string
	RetainHandling    uint8
	RetainAsPublished uint8
	NoLocal           uint8
	Qos               uint8
	ProtocolVersion   uint8
	Enabled           bool
	CreatedAt         int64
	SessionID         uint
	Shared            bool
	ShareName         string
}

type SubscriptionGroup map[string][]Subscription
