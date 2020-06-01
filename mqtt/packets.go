package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
)

func rangePackets(packets <-chan Packet, events chan<- Event, session *model.Session) {
	for p := range packets {
		// log.Println(p)
		switch p.PacketType() {
		case PACKET_TYPE_CONNECT:
			connectReq(p, events, session)
		case PACKET_TYPE_DISCONNECT:
			disconnectReq(p, events, session)
		case PACKET_TYPE_PUBLISH:
			publishReq(p, events, session)
		case PACKET_TYPE_PUBACK:
			pubackReq(p, events, session)
		case PACKET_TYPE_PUBREC:
			pubrecReq(p, events, session)
		case PACKET_TYPE_PUBREL:
			pubrelReq(p, events, session)
		case PACKET_TYPE_PUBCOMP:
			pubcompReq(p, events, session)
		case PACKET_TYPE_SUBSCRIBE:
			subscribeReq(p, events, session)
		case PACKET_TYPE_UNSUBSCRIBE:
			unsubscribeReq(p, events, session)
		case PACKET_TYPE_PINGREQ:
			pingReq(events, session)
		default:
			var event Event
			event.eventType = EVENT_PACKET_ERR
			log.Printf("Unknown packet type %d\n", p.PacketType())
			events <- event
		}
	}
}
