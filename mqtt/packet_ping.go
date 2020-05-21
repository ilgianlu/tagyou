package mqtt

import "github.com/ilgianlu/tagyou/model"

func pingReq(events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_PING
	event.clientId = session.ClientId
	event.session = session
	events <- event
}

func PingResp() Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_PINGRES) << 4
	p.remainingLength = 0
	return p
}
