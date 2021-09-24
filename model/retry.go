package model

import (
	"time"

	"github.com/ilgianlu/tagyou/conf"
)

// qos 1
const WAIT_FOR_PUB_ACK = 10

// qos 2
const WAIT_FOR_PUB_REC = 20
const WAIT_FOR_PUB_REL = 20
const WAIT_FOR_PUB_COMP = 21

type Retry struct {
	ID                 uint   `gorm:"primaryKey"`
	ClientId           string `gorm:"uniqueIndex:client_identifier_idx"`
	ApplicationMessage []byte
	PacketIdentifier   int `gorm:"uniqueIndex:client_identifier_idx"`
	Qos                uint8
	Dup                bool
	Retries            uint8
	AckStatus          uint8
	CreatedAt          int64
	SessionID          uint
	ReasonCode         uint8 `gorm:"-"`
}

func (r Retry) Expired() bool {
	return r.CreatedAt+int64(conf.RETRY_EXPIRATION) < time.Now().Unix()
}
