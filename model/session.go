package model

import (
	"net"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Session struct {
	ID              uint `gorm:"primary_key"`
	LastSeen        time.Time
	ExpiryInterval  uint32
	ClientId        string `gorm:"unique_index"`
	Connected       bool
	ProtocolVersion uint8     `gorm:"-"`
	ConnectFlags    uint8     `gorm:"-"`
	KeepAlive       int       `gorm:"-"`
	WillTopic       string    `gorm:"-"`
	WillDelay       time.Time `gorm:"-"`
	WillMessage     []byte    `gorm:"-"`
	Subscriptions   []Subscription
	Retries         []Retry
	Username        string   `gorm:"-"`
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

func (s Session) Expired() bool {
	return s.LastSeen.Add(time.Duration(s.ExpiryInterval) * time.Second).Before(time.Now())
}

func (s *Session) AfterDelete(tx *gorm.DB) (err error) {
	tx.Where("session_id = ?", s.ID).Delete(Subscription{})
	tx.Where("session_id = ?", s.ID).Delete(Retry{})
	return nil
}

func (s *Session) MergeSession(newSession Session) {
	s.ProtocolVersion = newSession.ProtocolVersion
	s.ConnectFlags = newSession.ConnectFlags
	s.KeepAlive = newSession.KeepAlive
	s.WillTopic = newSession.WillTopic
	s.WillDelay = newSession.WillDelay
	s.WillMessage = newSession.WillMessage
	s.Username = newSession.Username
	s.Password = newSession.Password
	s.ExpiryInterval = newSession.ExpiryInterval
	s.Conn = newSession.Conn
}

func CleanSession(db *gorm.DB, clientId string) {
	db.Where("client_id = ?", clientId).Delete(Session{})
}

func SessionExists(db *gorm.DB, clientId string) (Session, bool) {
	session := Session{}
	if db.Where("client_id = ?", clientId).First(&session).RecordNotFound() {
		return session, false
	} else {
		return session, true
	}
}

func DisconnectSession(db *gorm.DB, clientId string) {
	db.Model(&Session{}).Where("client_id = ?", clientId).Updates(map[string]interface{}{
		"Connected": false,
		"LastSeen":  time.Now(),
	})
}
