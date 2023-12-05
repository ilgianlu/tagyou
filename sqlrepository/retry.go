package sqlrepository

import (
	"github.com/ilgianlu/tagyou/model"
	"gorm.io/gorm"
)

type Retry struct {
	ID                 uint   `gorm:"primaryKey"`
	ClientId           string `gorm:"uniqueIndex:client_identifier_idx"`
	ApplicationMessage []byte
	PacketIdentifier   int `gorm:"uniqueIndex:client_identifier_idx"`
	Qos                uint8
	Dup                bool
	Retries            uint8
	AckStatus          uint8
	CreatedAt          int64
	SessionID          uint
	ReasonCode         uint8
}

type RetrySqlRepository struct {
	Db *gorm.DB
}

func (r RetrySqlRepository) InsertOne(retry model.Retry) error {
	return r.Db.Create(&retry).Error
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
