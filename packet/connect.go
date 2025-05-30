package packet

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

// ATTENTION: partial implementation only for testing
func Connect() Packet {
	var p Packet
	p.header = uint8(PACKET_TYPE_CONNECT) << 4
	return p
}

func (p *Packet) connectReq(session *model.RunningSession) int {
	// START VARIABLE HEADER
	i := 0
	pl := Read2BytesInt(p.remainingBytes, i)
	i = i + 2
	slog.Debug("protocolName", "bytes-read", pl, "protocol-name", string(p.remainingBytes[i:i+pl]))
	i = i + pl
	v := p.remainingBytes[i]
	slog.Debug("protocolVersion", "protocol-version", v)
	session.ProtocolVersion = v
	i++
	if int(v) < conf.MINIMUM_SUPPORTED_PROTOCOL {
		slog.Error("unsupported protocol version", "protocol-version", v)
		return UNSUPPORTED_PROTOCOL_VERSION
	}
	session.ConnectFlags = p.remainingBytes[i]
	i++
	ka := p.remainingBytes[i : i+2]
	session.KeepAlive = Read2BytesInt(ka, 0)
	slog.Debug("keepAlive", "keep-alive", Read2BytesInt(ka, 0))
	i = i + 2
	if session.ProtocolVersion >= conf.MQTT_V5 {
		pl, err := p.parseProperties(i)
		if err != 0 {
			slog.Error("err reading properties", "err", err)
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
				slog.Error("err reading properties", "err", err)
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
		slog.Debug("will topic with message", "will-topic", session.WillTopic, "will-message", session.WillMessage[:wml])
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
		slog.Debug("user logging in", "username", session.Username)
	}
	if session.ProtocolVersion >= conf.MQTT_V5 {
		session.ExpiryInterval = p.SessionExpiryInterval()
	}
	return 0
}
