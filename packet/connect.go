package packet

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/format"
	"github.com/ilgianlu/tagyou/model"
)

// Connect ATTENTION: partial implementation only for testing
func Connect() Packet {
	var p Packet
	p.header = header(uint8(PACKET_TYPE_CONNECT) << 4)
	return p
}

func (p *Packet) connectReq(session *model.RunningSession) int {
	// START VARIABLE HEADER
	i := 0
	pl, err := format.Read2BytesInt(p.remainingBytes, i)
	if err != nil {
		return 1
	}
	i = i + 2
	slog.Debug("[PACKET] protocolName", "bytes-read", pl, "protocol-name", string(p.remainingBytes[i:i+pl]))
	i = i + pl
	v := p.remainingBytes[i]
	slog.Debug("[PACKET] protocolVersion", "protocol-version", v)
	session.ProtocolVersion = v
	i++
	if int(v) < conf.MINIMUM_SUPPORTED_PROTOCOL {
		slog.Error("[PACKET] unsupported protocol version", "protocol-version", v)
		return UNSUPPORTED_PROTOCOL_VERSION
	}
	session.ConnectFlags = p.remainingBytes[i]
	i++
	ka := p.remainingBytes[i : i+2]
	session.KeepAlive, err = format.Read2BytesInt(ka, 0)
	if err != nil {
		return 1
	}
	slog.Debug("[PACKET] keepAlive", "keep-alive", session.KeepAlive)
	i = i + 2
	if session.ProtocolVersion >= conf.MQTT_V5 {
		pl, err := p.parseProperties(i)
		if err != 0 {
			slog.Error("[PACKET] err reading properties", "err", err)
			return err
		}
		i = i + pl
	}
	// END VARIABLE HEADER
	// START PAYLOAD
	p.payloadOffset = i
	cil, err := format.Read2BytesInt(p.remainingBytes, i)
	if err != nil {
		return 1
	}
	i = i + 2
	session.ClientId = string(p.remainingBytes[i : i+cil])
	// log.Printf("%d bytes, clientId %s\n", cil, event.clientId)
	i = i + cil
	if session.WillFlag() {
		if session.ProtocolVersion >= conf.MQTT_V5 {
			pl, err := p.parseWillProperties(i)
			if err != 0 {
				slog.Error("[PACKET] err reading properties", "err", err)
				return err
			}
			i = i + pl
		}
		// read will topic
		wtl, err := format.Read2BytesInt(p.remainingBytes, i)
		if err != nil {
			return 1
		}
		i = i + 2
		session.WillTopic = string(p.remainingBytes[i : i+wtl])
		i = i + wtl
		// will message
		wml, err := format.Read2BytesInt(p.remainingBytes, i)
		if err != nil {
			return 1
		}
		i = i + 2
		session.WillMessage = p.remainingBytes[i : i+wml]
		slog.Debug("[PACKET] will topic with message", "will-topic", session.WillTopic, "will-message", session.WillMessage[:wml])
		i = i + wml
	}
	if session.HaveUser() {
		// read username
		unl, err := format.Read2BytesInt(p.remainingBytes, i)
		if err != nil {
			return 1
		}
		i = i + 2
		session.Username = string(p.remainingBytes[i : i+unl])
		i = i + unl
		// read password
		pwdl, err := format.Read2BytesInt(p.remainingBytes, i)
		if err != nil {
			return 1
		}
		i = i + 2
		session.Password = string(p.remainingBytes[i : i+pwdl])
		slog.Debug("[PACKET] user logging in", "username", session.Username)
	}
	if session.ProtocolVersion >= conf.MQTT_V5 {
		session.ExpiryInterval = p.SessionExpiryInterval()
	}
	return 0
}
