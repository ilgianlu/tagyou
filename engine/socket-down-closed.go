package engine

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
)

func (s StandardEngine) OnSocketDownClosed(session *model.RunningSession) bool {
	slog.Debug("[MQTT] socket closed!", "client-id", session.GetClientId())
	if session.GetClientId() != "" {
		slog.Debug("[MQTT] client was connected!", "client-id", session.GetClientId(), "last-connect", session.LastConnect)
		sendWill(session)
		return true
	}
	return false
}
