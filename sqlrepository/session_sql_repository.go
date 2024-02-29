package sqlrepository

import (
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

type SessionSqlRepository struct {
	Db *dbaccess.Queries
}

func (sr SessionSqlRepository) PersistSession(running *model.RunningSession) (sessionId uint, err error) {
	running.Mu.RLock()
	defer running.Mu.RUnlock()
	sess := Session{
		LastSeen:        running.LastSeen,
		LastConnect:     running.LastConnect,
		ExpiryInterval:  running.ExpiryInterval,
		ClientId:        running.ClientId,
		Connected:       running.Connected,
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
	err := sr.Db.Where("client_id = ?", clientId).First(&session).Error
	return &session, err == nil
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
		return &session, err
	}

	return &session, nil
}

func (sr SessionSqlRepository) GetAll() []model.Session {
	sessions := []model.Session{}
	if err := sr.Db.Find(&sessions).Error; err != nil {
		return sessions
	}

	return sessions
}

func (sr SessionSqlRepository) Save(session *model.Session) {
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
