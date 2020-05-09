package model

import "time"

type Retain struct {
	Topic              string `gorm:"index:retain_topic_idx"`
	ApplicationMessage []byte
	CreatedAt          time.Time
}
