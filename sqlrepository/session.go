package sqlrepository

import (
	"context"
	"database/sql"
	"log/slog"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

type SessionSqlRepository struct {
	Db      *dbaccess.Queries
	SqlConn *sql.DB
}

func mapSession(session dbaccess.Session) model.Session {
	connectd := false
	if session.Connected.Int64 == 1 {
		connectd = true
	}
	return model.Session{
		ID:              session.ID,
		LastSeen:        session.LastSeen.Int64,
		LastConnect:     session.LastConnect.Int64,
		ExpiryInterval:  session.ExpiryInterval.Int64,
		ClientId:        session.ClientID.String,
		Connected:       connectd,
		ProtocolVersion: uint8(session.ProtocolVersion.Int64),
	}
}

func mappingSessions(sessions []dbaccess.Session) []model.Session {
	sesss := []model.Session{}
	for _, sess := range sessions {
		sesss = append(sesss, mapSession(sess))
	}
	return sesss
}

func (sr SessionSqlRepository) PersistSession(running *model.RunningSession) (int64, error) {
	var connectd int64 = 0
	if running.Connected {
		connectd = 1
	}
	sess := dbaccess.CreateSessionParams{
		LastSeen:        sql.NullInt64{Int64: running.LastSeen, Valid: true},
		LastConnect:     sql.NullInt64{Int64: running.LastConnect, Valid: true},
		ExpiryInterval:  sql.NullInt64{Int64: running.ExpiryInterval, Valid: true},
		ClientID:        sql.NullString{String: running.ClientId, Valid: true},
		Connected:       sql.NullInt64{Int64: connectd, Valid: true},
		ProtocolVersion: sql.NullInt64{Int64: int64(running.ProtocolVersion), Valid: true},
	}
	newSess, err := sr.Db.CreateSession(context.Background(), sess)
	return newSess.ID, err
}

func (sr SessionSqlRepository) UpdateSession(sessionId int64, running *model.RunningSession) (int64, error) {
	var connectd int64 = 0
	if running.Connected {
		connectd = 1
	}
	sess := dbaccess.UpdateSessionParams{
		ID:              sessionId,
		LastSeen:        sql.NullInt64{Int64: running.LastSeen, Valid: true},
		LastConnect:     sql.NullInt64{Int64: running.LastConnect, Valid: true},
		ExpiryInterval:  sql.NullInt64{Int64: running.ExpiryInterval, Valid: true},
		ClientID:        sql.NullString{String: running.ClientId, Valid: true},
		Connected:       sql.NullInt64{Int64: connectd, Valid: true},
		ProtocolVersion: sql.NullInt64{Int64: int64(running.ProtocolVersion), Valid: true},
	}
	newSess, err := sr.Db.UpdateSession(context.Background(), sess)
	return newSess.ID, err
}

func (sr SessionSqlRepository) CleanSession(clientId string) error {
	tx, err := sr.SqlConn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	qtx := sr.Db.WithTx(tx)
	sess, err := qtx.GetSessionByClientId(context.Background(), sql.NullString{String: clientId, Valid: true})
	if err != nil {
		return err
	}
	slog.Debug("[MQTT] cleaning session", "session-id", sess.ID)
	qtx.DeleteSubscriptionBySessionId(context.Background(), sql.NullInt64{Int64: sess.ID, Valid: true})
	qtx.DeleteRetryBySessionId(context.Background(), sql.NullInt64{Int64: sess.ID, Valid: true})
	qtx.DeleteSessionById(context.Background(), sess.ID)
	return tx.Commit()
}

func (sr SessionSqlRepository) SessionExists(clientId string) (model.Session, bool) {
	session, err := sr.Db.GetSessionByClientId(context.Background(), sql.NullString{String: clientId, Valid: true})
	return mapSession(session), err == nil
}

func (sr SessionSqlRepository) DisconnectSession(clientId string) {
	sr.Db.DisconnectSessionByClientId(context.Background(), dbaccess.DisconnectSessionByClientIdParams{
		ClientID: sql.NullString{String: clientId, Valid: true},
		LastSeen: sql.NullInt64{Int64: time.Now().Unix(), Valid: true},
	})
}

func (sr SessionSqlRepository) GetById(sessionId int64) (model.Session, error) {
	session, err := sr.Db.GetSessionById(context.Background(), sessionId)
	if err != nil {
		return mapSession(session), err
	}

	return mapSession(session), nil
}

func (sr SessionSqlRepository) GetAll() []model.Session {
	sessions, err := sr.Db.GetAllSessions(context.Background())

	if err != nil {
		return []model.Session{}
	}

	return mappingSessions(sessions)
}

func (sr SessionSqlRepository) IsOnline(sessionId int64) bool {
	session, err := sr.Db.GetSessionById(context.Background(), sessionId)
	if err != nil {
		return false
	} else {
		return session.Connected.Int64 == 1
	}
}
