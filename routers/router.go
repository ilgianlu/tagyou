package routers

import (
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
)

type Router interface {
	Send(clientId string, payload []byte)
	Forward(senderId string, topic string, packet *packet.Packet)
	SendRetain(protocolVersion uint8, subscription model.Subscription)
	AddDestination(clientId string, conn model.TagyouConn)
	RemoveDestination(clientId string)
	DestinationExists(clientId string) bool
}
