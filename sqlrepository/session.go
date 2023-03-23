package sqlrepository

import (
	"github.com/ilgianlu/tagyou/model"
	"gorm.io/gorm"
)

type Session struct {
	ID              uint `gorm:"primary_key"`
	LastSeen        int64
	LastConnect     int64
	ExpiryInterval  int64
	ClientId        string `gorm:"uniqueIndex:client_unique_session_idx"`
	Connected       bool
	ProtocolVersion uint8
	Subscriptions   []Subscription `json:"-"`
	Retries         []Retry        `json:"-"`
}

func (s *Session) BeforeDelete(tx *gorm.DB) (err error) {
	tx.Where("session_id = ?", s.ID).Delete(&Subscription{})
	tx.Where("session_id = ?", s.ID).Delete(&Retry{})
	return nil
}

func (s *Session) GetId() uint {
	return s.ID
}

func (s *Session) GetClientId() string {
	return s.ClientId
}

func (s *Session) GetProtocolVersion() uint8 {
	return s.ProtocolVersion
}

func (s *Session) Expired() bool {
	return model.SessionExpired(s.LastSeen, s.ExpiryInterval)
}
func (s *Session) GetLastSeen() int64 {
	return s.LastSeen
}

func (s *Session) GetLastConnect() int64 {
	return s.LastConnect
}
func (s *Session) GetExpiryInterval() int64 {
	return s.ExpiryInterval
}
func (s *Session) GetConnected() bool {
	return s.Connected
}
