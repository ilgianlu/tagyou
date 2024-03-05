package event

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
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
		if willPacket.Retain() {
			slog.Debug("[MQTT] will packet to retain")
			willPacket.Topic = session.WillTopic
			saveRetain(session, &willPacket)
		}
		router.Forward(session.WillTopic, &willPacket)
	}
}

func needWillSend(runningSession *model.RunningSession) bool {
	if session, ok := persistence.SessionRepository.SessionExists(runningSession.ClientId); ok {
		if session.GetLastConnect() > runningSession.LastConnect {
			// session persisted is newer then running memory session... device reconnected!
			// no need to send will
			slog.Debug("[MQTT] avoid sending will! (client reconnected)", "client-id", session.GetClientId())
			return false
		}
	}
	return true
}
