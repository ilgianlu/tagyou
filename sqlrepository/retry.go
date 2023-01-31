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
}

type RetrySqlRepository struct {
	Db *gorm.DB
}

func (r RetrySqlRepository) SaveOne(retry model.Retry) error {
	return r.Db.Save(&retry).Error
}
