package sqlrepository

import (
	"context"
	"database/sql"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

type RetrySqlRepository struct {
	Db *dbaccess.Queries
}

func (r RetrySqlRepository) InsertOne(retry model.Retry) error {
	dupl := int64(0)
	if retry.Dup {
		dupl = 1
	}
	params := dbaccess.CreateRetryParams{
		ClientID:           sql.NullString{String: retry.ClientId, Valid: true},
		ApplicationMessage: retry.ApplicationMessage,
		PacketIdentifier:   sql.NullInt64{Int64: int64(retry.PacketIdentifier), Valid: true},
		Qos:                sql.NullInt64{Int64: int64(retry.Qos), Valid: true},
		Dup:                sql.NullInt64{Int64: dupl, Valid: true},
		Retries:            sql.NullInt64{Int64: int64(retry.Retries), Valid: true},
		AckStatus:          sql.NullInt64{Int64: int64(retry.AckStatus), Valid: true},
		CreatedAt:          sql.NullInt64{Int64: retry.CreatedAt, Valid: true},
		SessionID:          sql.NullInt64{Int64: int64(retry.SessionID), Valid: true},
		ReasonCode:         sql.NullInt64{Int64: int64(retry.ReasonCode), Valid: true},
	}
	_, err := r.Db.CreateRetry(context.Background(), params)
	return err
}

func (r RetrySqlRepository) SaveOne(retry model.Retry) error {
	return r.Db.Save(&retry).Error
}

func (r RetrySqlRepository) Delete(retry model.Retry) error {
	return r.Db.Delete(&retry).Error
}

func (r RetrySqlRepository) FirstByClientIdPacketIdentifier(clientId string, packetIdentifier int) (model.Retry, error) {
	retry := model.Retry{
		ClientId:         clientId,
		PacketIdentifier: packetIdentifier,
	}
	err := r.Db.First(&retry).Error

	return retry, err
}

func (r RetrySqlRepository) FirstByClientIdPacketIdentifierReasonCode(clientId string, packetIdentifier int, reasonCode uint8) (model.Retry, error) {
	retry := model.Retry{
		ClientId:         clientId,
		PacketIdentifier: packetIdentifier,
		ReasonCode:       reasonCode,
	}
	err := r.Db.First(&retry).Error

	return retry, err
}
