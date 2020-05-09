package model

import "time"

type Subscription struct {
	ID                uint64
	ClientId          string `gorm:"index:clientid_idx"`
	Topic             string `gorm:"index:subscription_topic_idx"`
	RetainHandling    uint8
	RetainAsPublished uint8
	NoLocal           uint8
	QoS               uint8
	Enabled           bool
	CreatedAt         time.Time
	SessionID         uint64
}
