package sender

import (
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
)

type Sender interface {
	Send(clientId string, payload []byte)
	Forward(topic string, packet *packet.Packet)
	AddDestination(clientId string, conn model.TagyouConn)
	RemoveDestination(clientId string)
	DestinationExists(clientId string) bool
}
