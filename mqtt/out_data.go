package mqtt

import "github.com/ilgianlu/tagyou/packet"

type OutData struct {
	clientId string
	packet   packet.Packet
}
