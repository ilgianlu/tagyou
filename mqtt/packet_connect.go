package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func connectReq(p Packet, events chan<- Event, session *model.Session) {
	var event Event
	event.eventType = EVENT_CONNECT
	// START VARIABLE HEADER
	i := 0
	pl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	// log.Printf("%d bytes, protocolName %s\n", pl, string(p.remainingBytes[i:i+pl]))
	i = i + pl
	v := p.remainingBytes[i]
	// log.Println("protocolVersion", v)
	session.ProtocolVersion = v
	i++
	if int(v) < conf.MINIMUM_SUPPORTED_PROTOCOL {
		log.Println("unsupported protocol version err", v)
		event.err = UNSUPPORTED_PROTOCOL_VERSION
		events <- event
		return
	}
	session.ConnectFlags = p.remainingBytes[i]
	i++
	ka := p.remainingBytes[i : i+2]
	session.KeepAlive = Read2BytesInt(ka, 0)
	// log.Println("keepAlive", Read2BytesInt(ka, 0))
	i = i + 2
	if session.ProtocolVersion >= MQTT_V5 {
		pl, pp, err := p.parseProperties(i)
		if err != 0 {
			log.Println("err reading properties", err)
			event.err = uint8(err)
			events <- event
			return
		}
		p.propertiesPos = pp
		i = i + pl
	}
	// END VARIABLE HEADER
	// START PAYLOAD
	cil := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	event.clientId = string(p.remainingBytes[i : i+cil])
	session.ClientId = event.clientId
	log.Printf("%d bytes, clientId %s\n", cil, event.clientId)
	i = i + cil
	if session.WillFlag() {
		if session.ProtocolVersion >= MQTT_V5 {
			pl, pp, err := p.parseProperties(i)
			if err != 0 {
				log.Println("err reading properties", err)
				event.err = uint8(err)
				events <- event
				return
			}
			p.willPropertiesLength = pl
			p.willPropertiesPos = pp
			i = i + pl
		}
		// read will topic
		wtl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		session.WillTopic = string(p.remainingBytes[i : i+wtl])
		i = i + wtl
		// will message
		wml := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		session.WillMessage = p.remainingBytes[i : i+wml]
		log.Printf("will topic \"%s\"\nwith message \"%s\"\n", session.WillTopic, session.WillMessage)
		i = i + wml
	}
	if session.HaveUser() {
		// read username
		unl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		session.Username = string(p.remainingBytes[i : i+unl])
		i = i + unl
		// read password
		pwdl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		session.Password = string(p.remainingBytes[i : i+pwdl])
		log.Printf("username \"%s\"\nlogging with password \"%s\"\n", session.Username, session.Password)
	}
	event.session = session
	events <- event
}
