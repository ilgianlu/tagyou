package event

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
)

func OnClientDisconnect(session *model.RunningSession, clientId string) {
	session.SetConnected(false)
	if session.Router.DestinationExists(clientId) {
		needDisconnection := needDisconnection(session)
		if !needDisconnection {
			return
		}
		session.Router.RemoveDestination(clientId)
		persistence.SessionRepository.DisconnectSession(clientId)
	}
}

func needDisconnection(runningSession *model.RunningSession) bool {
	if session, ok := persistence.SessionRepository.SessionExists(runningSession.ClientId); ok {
		if session.GetLastConnect() > runningSession.LastConnect {
			// session persisted is newer then running memory session... device reconnected!
			// no need to send will
			slog.Debug("[MQTT] avoid disconnect! (client reconnected)", "client-id", session.GetClientId())
			return false
		}
	}
	return true
}
