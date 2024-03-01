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

func mappingRetry(retry dbaccess.Retry) model.Retry {
	dupl := false
	if retry.Dup.Int64 == 1 {
		dupl = true
	}
	return model.Retry{
		ID:                 retry.ID,
		ClientId:           retry.ClientID.String,
		ApplicationMessage: retry.ApplicationMessage,
		PacketIdentifier:   int(retry.PacketIdentifier.Int64),
		Qos:                uint8(retry.Qos.Int64),
		Dup:                dupl,
		Retries:            uint8(retry.Retries.Int64),
		AckStatus:          uint8(retry.AckStatus.Int64),
		CreatedAt:          retry.CreatedAt.Int64,
		SessionID:          retry.SessionID.Int64,
		ReasonCode:         uint8(retry.ReasonCode.Int64),
	}
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

func (r RetrySqlRepository) UpdateAckStatus(id uint, ackStatus uint) error {
	return r.Db.UpdateRetryAckStatus(context.Background(), dbaccess.UpdateRetryAckStatusParams{
		AckStatus: sql.NullInt64{Int64: int64(ackStatus), Valid: true},
		ID:        int64(id),
	})
}

func (r RetrySqlRepository) Delete(retry model.Retry) error {
	return r.Db.DeleteRetryById(context.Background(), int64(retry.ID))
}

func (r RetrySqlRepository) FirstByClientIdPacketIdentifier(clientId string, packetIdentifier int) (model.Retry, error) {
	params := dbaccess.GetRetryByClientIdPacketIdentifierParams{
		ClientID:         sql.NullString{String: clientId, Valid: true},
		PacketIdentifier: sql.NullInt64{Int64: int64(packetIdentifier), Valid: true},
	}
	retry, err := r.Db.GetRetryByClientIdPacketIdentifier(context.Background(), params)
	if err != nil {
		return model.Retry{}, err
	}
	return mappingRetry(retry), nil
}

func (r RetrySqlRepository) FirstByClientIdPacketIdentifierReasonCode(clientId string, packetIdentifier int, reasonCode uint8) (model.Retry, error) {
	params := dbaccess.GetRetryByClientIdPacketIdentifierReasonCodeParams{
		ClientID:         sql.NullString{String: clientId, Valid: true},
		PacketIdentifier: sql.NullInt64{Int64: int64(packetIdentifier), Valid: true},
		ReasonCode:       sql.NullInt64{Int64: int64(reasonCode), Valid: true},
	}
	retry, err := r.Db.GetRetryByClientIdPacketIdentifierReasonCode(context.Background(), params)
	if err != nil {
		return model.Retry{}, err
	}
	return mappingRetry(retry), err
}
