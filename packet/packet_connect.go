package packet

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func (p *Packet) connectReq(session *model.RunningSession) int {
	session.Mu.Lock()
	defer session.Mu.Unlock()
	p.Event = EVENT_CONNECT
	// START VARIABLE HEADER
	i := 0
	pl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	log.Debug().Msgf("%d bytes, protocolName %s\n", pl, string(p.remainingBytes[i:i+pl]))
	i = i + pl
	v := p.remainingBytes[i]
	log.Debug().Msgf("protocolVersion %d", v)
	session.ProtocolVersion = v
	i++
	if int(v) < conf.MINIMUM_SUPPORTED_PROTOCOL {
		log.Error().Msgf("unsupported protocol version err %d", v)
		return UNSUPPORTED_PROTOCOL_VERSION
	}
	session.ConnectFlags = p.remainingBytes[i]
	i++
	ka := p.remainingBytes[i : i+2]
	session.KeepAlive = Read2BytesInt(ka, 0)
	log.Debug().Msgf("keepAlive %d", Read2BytesInt(ka, 0))
	i = i + 2
	if session.ProtocolVersion >= conf.MQTT_V5 {
		pl, err := p.parseProperties(i)
		if err != 0 {
			log.Error().Msgf("err reading properties %d", err)
			return err
		}
		i = i + pl
	}
	// END VARIABLE HEADER
	// START PAYLOAD
	p.payloadOffset = i
	cil := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	session.ClientId = string(p.remainingBytes[i : i+cil])
	// log.Printf("%d bytes, clientId %s\n", cil, event.clientId)
	i = i + cil
	if session.WillFlag() {
		if session.ProtocolVersion >= conf.MQTT_V5 {
			pl, err := p.parseWillProperties(i)
			if err != 0 {
				log.Error().Msgf("err reading properties %d", err)
				return err
			}
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
		log.Debug().Msgf("will topic %s with message %s", session.WillTopic, session.WillMessage[:20])
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
		log.Debug().Msgf("username \"%s\"\nlogging with password \"%s\"\n", session.Username, session.Password)
	}
	if session.ProtocolVersion >= conf.MQTT_V5 {
		session.ExpiryInterval = p.SessionExpiryInterval()
	}
	return 0
}
