package model

import "time"

type Retain struct {
	Topic              string `gorm:"primary_key;auto_increment:false"`
	ApplicationMessage []byte
	CreatedAt          time.Time
}
