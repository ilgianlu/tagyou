package event

import (
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sender"
	"github.com/rs/zerolog/log"
)

func sendWill(sender sender.Sender, p *packet.Packet) {
	p.Session.Mu.RLock()
	defer p.Session.Mu.RUnlock()
	if p.Session.WillTopic != "" {
		needWillSend := needWillSend(p)
		if !needWillSend {
			return
		}
		willPacket := packet.Publish(p.Session.ProtocolVersion, p.Session.WillQoS(), p.Session.WillRetain(), p.Session.WillTopic, packet.NewPacketIdentifier(), p.Session.WillMessage)
		sender.Forward(p.Session.WillTopic, &willPacket)
	}
}

func needWillSend(p *packet.Packet) bool {
	if session, ok := persistence.SessionRepository.SessionExists(p.Session.ClientId); ok {
		log.Debug().Msgf("[MQTT] (%s) Persisted session LastConnect %d running session %d", p.Session.ClientId, session.LastConnect, p.Session.LastConnect)
		if session.LastConnect > p.Session.LastConnect {
			// session persisted is newer then running memory session... device reconnected!
			// no need to send will
			log.Debug().Msgf("[MQTT] (%s) avoid sending will! (device reconnected)", p.Session.ClientId)
			return false
		}
	}
	return true
}
