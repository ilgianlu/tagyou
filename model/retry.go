package model

import (
	"time"
)

// qos 1
const WAIT_FOR_PUB_ACK = 10

// qos 2
const WAIT_FOR_PUB_REC = 20
const WAIT_FOR_PUB_REL = 20
const WAIT_FOR_PUB_COMP = 21

type Retry struct {
	ClientId           string `gorm:"primary_key;auto_increment:false"`
	ApplicationMessage []byte
	PacketIdentifier   int `gorm:"primary_key;auto_increment:false"`
	Qos                uint8
	Dup                bool
	Retries            uint8
	AckStatus          uint8
	CreatedAt          time.Time
	SessionID          uint
	ReasonCode         uint8 `gorm:"-"`
}
