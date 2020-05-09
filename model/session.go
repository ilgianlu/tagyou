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

func (s Session) ReservedBit() bool {
	return (s.ConnectFlags & 0x01) == 0
}

func (s Session) CleanStart() bool {
	return (s.ConnectFlags & 0x02) > 0
}

func (s Session) WillFlag() bool {
	return (s.ConnectFlags & 0x04) > 0
}

func (s Session) WillQoS() uint8 {
	return (s.ConnectFlags & 0x18 >> 3)
}

func (s Session) WillRetain() bool {
	return (s.ConnectFlags & 0x20) > 0
}

func (s Session) HavePass() bool {
	return (s.ConnectFlags & 0x40) > 0
}

func (s Session) HaveUser() bool {
	return (s.ConnectFlags & 0x80) > 0
}
