package model

import (
	"net"
	"strings"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"gorm.io/gorm"
)

type Session struct {
	ID              uint `gorm:"primary_key"`
	LastSeen        int64
	ExpiryInterval  int64
	ClientId        string `gorm:"uniqueIndex:client_unique_session_idx"`
	Connected       bool
	ProtocolVersion uint8
	ConnectFlags    uint8          `gorm:"-" json:"-"`
	KeepAlive       int            `gorm:"-" json:"-"`
	WillTopic       string         `gorm:"-" json:"-"`
	WillDelay       int64          `gorm:"-" json:"-"`
	WillMessage     []byte         `gorm:"-" json:"-"`
	Subscriptions   []Subscription `json:"-"`
	Retries         []Retry        `json:"-"`
	Username        string         `gorm:"-" json:"-"`
	Password        string         `gorm:"-" json:"-"`
	SubscribeAcl    string         `gorm:"-" json:"-"`
	PublishAcl      string         `gorm:"-" json:"-"`
	Conn            net.Conn       `gorm:"-" json:"-"`
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
	return s.LastSeen+s.ExpiryInterval < time.Now().Unix()
}

func (s Session) FromLocalhost() bool {
	return strings.Index(s.Conn.RemoteAddr().String(), conf.LOCALHOST) == 0
}

func (s *Session) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Where("session_id = ?", s.ID).Delete(&Subscription{})
	tx.Where("session_id = ?", s.ID).Delete(&Retry{})
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

func CleanSession(db *gorm.DB, clientId string) error {
	sess := Session{}
	if err := db.Where("client_id = ?", clientId).First(&sess).Error; err != nil {
		return err
	}
	return db.Delete(&sess).Error
}

func SessionExists(db *gorm.DB, clientId string) (Session, bool) {
	session := Session{}
	if err := db.Where("client_id = ?", clientId).First(&session).Error; err != nil {
		return session, false
	} else {
		return session, true
	}
}

func DisconnectSession(db *gorm.DB, clientId string) {
	db.Model(&Session{}).Where("client_id = ?", clientId).Updates(map[string]interface{}{
		"Connected": false,
		"LastSeen":  time.Now().Unix(),
	})
}
