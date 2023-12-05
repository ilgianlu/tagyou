package repository

import "github.com/ilgianlu/tagyou/model"

type RetryRepository interface {
	InsertOne(retry model.Retry) error
	SaveOne(retry model.Retry) error
	Delete(retry model.Retry) error
	FirstByClientIdPacketIdentifier(clientId string, packetIdentifier int) (model.Retry, error)
	FirstByClientIdPacketIdentifierReasonCode(clientId string, packetIdentifier int, reasonCode uint8) (model.Retry, error)
}
