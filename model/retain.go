package model

import "time"

type Retain struct {
	Topic              string `gorm:"primaryKey"`
	ApplicationMessage []byte
	CreatedAt          time.Time
}
