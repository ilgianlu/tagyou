package event

import (
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

func sendWill(db *gorm.DB, p *packet.Packet, outQueue chan<- out.OutData) {
	p.Session.Mu.RLock()
	defer p.Session.Mu.RUnlock()
	if p.Session.WillTopic != "" {
		needWillSend := needWillSend(db, p)
		if !needWillSend {
			return
		}
		willPacket := packet.Publish(p.Session.ProtocolVersion, p.Session.WillQoS(), p.Session.WillRetain(), p.Session.WillTopic, packet.NewPacketIdentifier(), p.Session.WillMessage)
		sendForward(db, p.Session.WillTopic, &willPacket, outQueue)
	}
}

func needWillSend(db *gorm.DB, p *packet.Packet) bool {
	if session, ok := model.SessionExists(db, p.Session.ClientId); ok {
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
