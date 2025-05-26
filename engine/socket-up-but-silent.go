package engine

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
)

func (s StandardEngine) OnSocketUpButSilent(session *model.RunningSession) bool {
	slog.Debug("[MQTT] keepalive not respected!", "keep-alive", session.KeepAlive*2)
	if session.GetClientId() != "" {
		slog.Debug("[MQTT] will due to keepalive not respected!", "client-id", session.GetClientId(), "last-connect", session.LastConnect)
		sendWill(session)
		return true
	}
	return false
}
