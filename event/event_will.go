package event

import (
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
	"github.com/rs/zerolog/log"
)

func sendWill(router routers.Router, session *model.RunningSession) {
	session.Mu.RLock()
	defer session.Mu.RUnlock()
	if session.WillTopic != "" {
		needWillSend := needWillSend(session)
		if !needWillSend {
			return
		}
		willPacket := packet.Publish(session.ProtocolVersion, session.WillQoS(), session.WillRetain(), session.WillTopic, packet.NewPacketIdentifier(), session.WillMessage)
		router.Forward(session.WillTopic, &willPacket)
	}
}

func needWillSend(runningSession *model.RunningSession) bool {
	if session, ok := persistence.SessionRepository.SessionExists(runningSession.ClientId); ok {
		log.Debug().Msgf("[MQTT] (%s) Persisted session LastConnect %d running session %d", session.GetClientId(), session.GetLastConnect(), runningSession.LastConnect)
		if session.GetLastConnect() > runningSession.LastConnect {
			// session persisted is newer then running memory session... device reconnected!
			// no need to send will
			log.Debug().Msgf("[MQTT] (%s) avoid sending will! (device reconnected)", session.GetClientId())
			return false
		}
	}
	return true
}
