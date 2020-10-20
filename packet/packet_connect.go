package packet

import (
	"log"

	"github.com/ilgianlu/tagyou/conf"
)

func (p *Packet) connectReq() int {
	p.Event = EVENT_CONNECT
	// START VARIABLE HEADER
	i := 0
	pl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	// log.Printf("%d bytes, protocolName %s\n", pl, string(p.remainingBytes[i:i+pl]))
	i = i + pl
	v := p.remainingBytes[i]
	// log.Println("protocolVersion", v)
	p.Session.ProtocolVersion = v
	i++
	if int(v) < conf.MINIMUM_SUPPORTED_PROTOCOL {
		log.Println("unsupported protocol version err", v)
		return UNSUPPORTED_PROTOCOL_VERSION
	}
	p.Session.ConnectFlags = p.remainingBytes[i]
	i++
	ka := p.remainingBytes[i : i+2]
	p.Session.KeepAlive = Read2BytesInt(ka, 0)
	// log.Println("keepAlive", Read2BytesInt(ka, 0))
	i = i + 2
	if p.Session.ProtocolVersion >= conf.MQTT_V5 {
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
	p.Session.ClientId = string(p.remainingBytes[i : i+cil])
	// log.Printf("%d bytes, clientId %s\n", cil, event.clientId)
	i = i + cil
	if p.Session.WillFlag() {
		if p.Session.ProtocolVersion >= conf.MQTT_V5 {
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
		p.Session.WillTopic = string(p.remainingBytes[i : i+wtl])
		i = i + wtl
		// will message
		wml := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		p.Session.WillMessage = p.remainingBytes[i : i+wml]
		// log.Printf("will topic \"%s\"\nwith message \"%s\"\n", p.Session.WillTopic, p.Session.WillMessage)
		i = i + wml
	}
	if p.Session.HaveUser() {
		// read username
		unl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		p.Session.Username = string(p.remainingBytes[i : i+unl])
		i = i + unl
		// read password
		pwdl := Read2BytesInt(p.remainingBytes, i)
		i = i + 2
		p.Session.Password = string(p.remainingBytes[i : i+pwdl])
		log.Printf("username \"%s\"\nlogging with password \"%s\"\n", p.Session.Username, p.Session.Password)
	}
	if p.Session.ProtocolVersion >= conf.MQTT_V5 {
		p.Session.ExpiryInterval = p.SessionExpiryInterval()
	}
	return 0
}
