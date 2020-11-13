package model

import "time"

type Subscription struct {
	ID                uint   `gorm:"primary_key"`
	ClientId          string `gorm:"uniqueIndex:sub_pars_idx"`
	Topic             string `gorm:"uniqueIndex:sub_pars_idx"`
	RetainHandling    uint8
	RetainAsPublished uint8
	NoLocal           uint8
	Qos               uint8
	ProtocolVersion   uint8
	Enabled           bool
	CreatedAt         time.Time
	SessionID         uint
	Shared            bool   `gorm:"default:false"`
	ShareName         string `gorm:"uniqueIndex:sub_pars_idx"`
}
