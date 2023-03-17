package sender

import "github.com/ilgianlu/tagyou/packet"

type Sender interface {
	Send(clientId string, payload []byte)
	Forward(topic string, packet packet.Packet)
}
