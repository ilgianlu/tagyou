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

func MapSession(session Session) model.Session {
	mSession := model.Session{
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

func MapSessions(sessions []Session) []model.Session {
	mSessions := []model.Session{}
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

func (sr SessionSqlRepository) SessionExists(clientId string) (model.Session, bool) {
	session := Session{}
	if err := sr.Db.Where("client_id = ?", clientId).First(&session).Error; err != nil {
		return model.Session{}, false
	} else {
		mSession := model.Session{
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

func (sr SessionSqlRepository) GetById(sessionId uint) (model.Session, error) {
	var session Session
	if err := sr.Db.Where("id = ?", sessionId).First(&session).Error; err != nil {
		return model.Session{}, err
	}

	mSession := model.Session{
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

func (sr SessionSqlRepository) GetAll() []model.Session {
	sessions := []Session{}
	if err := sr.Db.Find(&sessions).Error; err != nil {
		return []model.Session{}
	}

	return MapSessions(sessions)
}

func (sr SessionSqlRepository) Save(session *model.Session) {
	sr.Db.Save(&session)
}
