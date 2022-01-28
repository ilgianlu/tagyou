package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID              uint `gorm:"primary_key"`
	LastSeen        int64
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

func (s Session) Expired() bool {
	return s.LastSeen+s.ExpiryInterval < time.Now().Unix()
}

func PersistSession(db *gorm.DB, running *RunningSession, connected bool) (sessionId uint, err error) {
	sess := Session{
		LastSeen:        running.LastSeen,
		ExpiryInterval:  running.ExpiryInterval,
		ClientId:        running.ClientId,
		Connected:       connected,
		ProtocolVersion: running.ProtocolVersion,
	}
	saveErr := db.Save(&sess).Error
	return sess.ID, saveErr
}

func (s *Session) UpdateFromRunning(running *RunningSession) {
	running.Mu.RLock()
	s.ProtocolVersion = running.ProtocolVersion
	s.ExpiryInterval = running.ExpiryInterval
	s.Connected = true
	running.Mu.RUnlock()
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
