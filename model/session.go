package model

import (
	"net"
	"time"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Session struct {
	ID              uint64
	ExpireAt        time.Time
	ClientId        string `gorm:"unique"`
	Connected       bool
	ProtocolVersion uint8
	ConnectFlags    uint8
	KeepAlive       int
	WillTopic       string
	WillDelay       time.Time
	WillMessage     []byte
	Subscriptions   []Subscription
	Retries         []Retry
	Conn            net.Conn `gorm:"-"`
}

func (s Session) reservedBit() bool {
	return (s.ConnectFlags & 0x01) == 0
}

func (s Session) cleanStart() bool {
	return (s.ConnectFlags & 0x02) > 0
}

func (s Session) willFlag() bool {
	return (s.ConnectFlags & 0x04) > 0
}
