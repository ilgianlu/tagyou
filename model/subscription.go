package model

import "time"

type Subscription struct {
	ID                uint   `gorm:"primary_key"`
	ClientId          string `gorm:"unique_index:sub_pars_idx"`
	Topic             string `gorm:"unique_index:sub_pars_idx"`
	RetainHandling    uint8
	RetainAsPublished uint8
	NoLocal           uint8
	QoS               uint8
	Enabled           bool
	CreatedAt         time.Time
	SessionID         uint
}
