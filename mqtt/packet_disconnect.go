package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
)

func disconnectReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_DISCONNECT
	event.clientId = session.ClientId
	event.session = session
	i := 0
	p.reasonCode = p.remainingBytes[i]
	if session.ProtocolVersion >= MQTT_V5 {
		pl, pp, err := p.readProperties(i)
		if err != 0 {
			log.Println("err reading properties", err)
			event.err = uint8(err)
			events <- event
			return
		}
		p.propertiesLength = pl
		p.propertiesPos = pp
	}
	events <- event
}
