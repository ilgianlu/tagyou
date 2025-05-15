package model

type Router interface {
	Send(clientId string, payload []byte)
	Forward(senderId string, topic string, packet Packet)
	SendRetain(protocolVersion uint8, subscription Subscription)
	AddDestination(clientId string, conn TagyouConn)
	RemoveDestination(clientId string)
	DestinationExists(clientId string) bool
}
