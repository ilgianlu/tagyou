package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/conf"
)

func (p *Packet) connectReq() int {
	p.event = EVENT_CONNECT
	// START VARIABLE HEADER
	i := 0
	pl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	// log.Printf("%d bytes, protocolName %s\n", pl, string(p.remainingBytes[i:i+pl]))
	i = i + pl
	v := p.remainingBytes[i]
	// log.Println("protocolVersion", v)
	p.session.ProtocolVersion = v
	i++
	if int(v) < conf.MINIMUM_SUPPORTED_PROTOCOL {
		log.Println("unsupported protocol version err", v)
		return UNSUPPORTED_PROTOCOL_VERSION
	}
	p.session.ConnectFlags = p.remainingBytes[i]
	i++
	ka := p.remainingBytes[i : i+2]
	p.session.KeepAlive = Read2BytesInt(ka, 0)
	// log.Println("keepAlive", Read2BytesInt(ka, 0))
	i = i + 2
	if p.session.ProtocolVersion >= MQTT_V5 {
		pl, err := p.parseProperties(i)
		if err != 0 {
			log.Println("err reading properties", err)
			return err
		}
		i = i + pl
	}
	// END VARIABLE HEADER
	// START PAYLOAD
	p.payloadOffset = i
	cil := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	p.session.ClientId = string(p.remainingBytes[i : i+cil])
	// log.Printf("%d bytes, clientId %s\n", cil, event.clientId)
	i = i + cil
	if p.session.WillFlag() {
		if p.session.ProtocolVersion >= MQTT_V5 {
			pl, err := p.parseWillProperties(i)
			if err != 0 {
				log.Println("err reading properties", err)
				return err
			}
			i = i + pl
		}
		// read will topic
		wtl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		p.session.WillTopic = string(p.remainingBytes[i : i+wtl])
		i = i + wtl
		// will message
		wml := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		p.session.WillMessage = p.remainingBytes[i : i+wml]
		log.Printf("will topic \"%s\"\nwith message \"%s\"\n", p.session.WillTopic, p.session.WillMessage)
		i = i + wml
	}
	if p.session.HaveUser() {
		// read username
		unl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		p.session.Username = string(p.remainingBytes[i : i+unl])
		i = i + unl
		// read password
		pwdl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		p.session.Password = string(p.remainingBytes[i : i+pwdl])
		log.Printf("username \"%s\"\nlogging with password \"%s\"\n", p.session.Username, p.session.Password)
	}
	if p.session.ProtocolVersion >= MQTT_V5 {
		p.session.ExpiryInterval = p.SessionExpiryInterval()
	}
	return 0
}
