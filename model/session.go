package model

import (
	"net"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Session struct {
	ID              uint `gorm:"primary_key"`
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
	Username        string
	Password        string   `gorm:"-"`
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

func (s *Session) AfterDelete(tx *gorm.DB) (err error) {
	tx.Where("session_id = ?", s.ID).Delete(Subscription{})
	tx.Where("session_id = ?", s.ID).Delete(Retry{})
	return nil
}
