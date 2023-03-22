package sqlrepository

import (
	"time"

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

type SessionSqlRepository struct {
	Db *gorm.DB
}

func MapSession(session Session) model.RunningSession {
	mSession := model.RunningSession{
		ID:              session.ID,
		LastSeen:        session.LastSeen,
		LastConnect:     session.LastConnect,
		ExpiryInterval:  session.ExpiryInterval,
		ClientId:        session.ClientId,
		Connected:       session.Connected,
		ProtocolVersion: session.ProtocolVersion,
	}
	return mSession
}

func MapSessions(sessions []Session) []model.RunningSession {
	mSessions := []model.RunningSession{}
	for _, s := range sessions {
		mSessions = append(mSessions, MapSession(s))
	}
	return mSessions
}

func (sr SessionSqlRepository) PersistSession(running *model.RunningSession, connected bool) (sessionId uint, err error) {
	running.Mu.RLock()
	defer running.Mu.RUnlock()
	sess := Session{
		LastSeen:        running.LastSeen,
		LastConnect:     running.LastConnect,
		ExpiryInterval:  running.ExpiryInterval,
		ClientId:        running.ClientId,
		Connected:       connected,
		ProtocolVersion: running.ProtocolVersion,
	}
	saveErr := sr.Db.Save(&sess).Error
	return sess.ID, saveErr
}

func (sr SessionSqlRepository) CleanSession(clientId string) error {
	sess := Session{}
	if err := sr.Db.Where("client_id = ?", clientId).First(&sess).Error; err != nil {
		return err
	}
	return sr.Db.Delete(&sess).Error
}

func (sr SessionSqlRepository) SessionExists(clientId string) (model.RunningSession, bool) {
	session := Session{}
	if err := sr.Db.Where("client_id = ?", clientId).First(&session).Error; err != nil {
		return model.RunningSession{}, false
	} else {
		mSession := model.RunningSession{
			ID:              session.ID,
			LastSeen:        session.LastSeen,
			LastConnect:     session.LastConnect,
			ExpiryInterval:  session.ExpiryInterval,
			ClientId:        session.ClientId,
			Connected:       session.Connected,
			ProtocolVersion: session.ProtocolVersion,
		}
		return mSession, true
	}
}

func (sr SessionSqlRepository) DisconnectSession(clientId string) {
	sr.Db.Model(&Session{}).Where("client_id = ?", clientId).Updates(map[string]interface{}{
		"Connected": false,
		"LastSeen":  time.Now().Unix(),
	})
}

func (sr SessionSqlRepository) GetById(sessionId uint) (model.RunningSession, error) {
	var session Session
	if err := sr.Db.Where("id = ?", sessionId).First(&session).Error; err != nil {
		return model.RunningSession{}, err
	}

	mSession := model.RunningSession{
		ID:              session.ID,
		LastSeen:        session.LastSeen,
		LastConnect:     session.LastConnect,
		ExpiryInterval:  session.ExpiryInterval,
		ClientId:        session.ClientId,
		Connected:       session.Connected,
		ProtocolVersion: session.ProtocolVersion,
	}
	return mSession, nil
}

func (sr SessionSqlRepository) GetAll() []model.RunningSession {
	sessions := []Session{}
	if err := sr.Db.Find(&sessions).Error; err != nil {
		return []model.RunningSession{}
	}

	return MapSessions(sessions)
}

func (sr SessionSqlRepository) Save(session *model.RunningSession) {
	sr.Db.Save(&session)
}

func (sr SessionSqlRepository) IsOnline(sessionId uint) bool {
	session := model.RunningSession{}
	if err := sr.Db.Where("id = ?", sessionId).First(&session).Error; err != nil {
		return false
	} else {
		return session.Connected
	}
}
