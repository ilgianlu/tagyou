package engine

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/password"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func (s StandardEngine) OnConnect(session *model.RunningSession) {
	clientId := session.GetClientId()
	if conf.FORBID_ANONYMOUS_LOGIN && !session.FromLocalhost() {
		if !doAuth(session) {
			return
		}
	}
	taken := checkConnectionTakeOver(session)
	if taken {
		slog.Debug("[MQTT] client reconnecting", "client-id", clientId)
	}
	session.Router.AddDestination(clientId, session.GetConn())

	startSession(session)

	connack := packet.Connack(false, packet.CONNECT_OK, session.GetProtocolVersion())
	session.Router.Send(clientId, connack.ToByteSlice())
}

func doAuth(session *model.RunningSession) bool {
	clientId := session.ClientId
	username := session.Username
	sessionPassword := session.Password
	ok, pubAcl, subAcl := checkAuth(clientId, username, sessionPassword)
	if !ok {
		slog.Debug("[MQTT] wrong connect credentials")
		return false
	}
	session.ApplyAcl(pubAcl, subAcl)
	slog.Debug("[MQTT] auth ok, imported acls", "pub-acl", pubAcl, "sub-acl", subAcl)
	return true
}

func checkAuth(clientId string, username string, sessionPassword string) (bool, string, string) {
	client, err := persistence.ClientRepository.GetByClientIdUsername(clientId, username)
	if err != nil {
		slog.Debug("[MQTT] could not find user")
		return false, "", ""
	}

	if err := password.CheckPassword(client.Password, []byte(sessionPassword)); err != nil {
		slog.Debug("[MQTT] wrong password")
		return false, "", ""
	}
	return true, client.PublishAcl, client.SubscribeAcl
}

func checkConnectionTakeOver(session *model.RunningSession) bool {
	clientId := session.ClientId
	protocolVersion := session.ProtocolVersion
	if !session.Router.DestinationExists(clientId) {
		return false
	}

	pkt := packet.Connack(false, packet.SESSION_TAKEN_OVER, protocolVersion)
	session.Router.Send(clientId, pkt.ToByteSlice())

	session.Router.RemoveDestination(clientId)
	return true
}

func startSession(session *model.RunningSession) {
	clientId := session.GetClientId()
	session.Router = routers.ByClientId(clientId, session.Router.GetConns())
	if prevSession, ok := persistence.SessionRepository.SessionExists(clientId); ok {
		slog.Debug("[MQTT] check existing session", "last-seen", prevSession.GetLastSeen(), "clean-start", session.CleanStart(), "expired", prevSession.Expired(), "new-protocol-version", session.GetProtocolVersion(), "prev-protocol-version", prevSession.GetProtocolVersion())
		if session.CleanStart() || prevSession.Expired() || session.GetProtocolVersion() != prevSession.GetProtocolVersion() {
			slog.Debug("[MQTT] Cleaning previous session: Invalid or to clean", "client-id", clientId)
			if err := persistence.SessionRepository.CleanSession(clientId); err != nil {
				slog.Error("[MQTT] error removing previous session", "client-id", clientId, "err", err)
			}
			session.SetConnected(true)
			if id, err := persistence.SessionRepository.PersistSession(session); err != nil {
				slog.Error("[MQTT] error persisting clean session", "client-id", clientId, "err", err)
			} else {
				session.ApplySessionId(id)
			}
		} else {
			slog.Debug("Updating previous session from running", "client-id", clientId)
			session.ApplySessionId(prevSession.GetId())
			session.SetConnected(true)
			persistence.SessionRepository.UpdateSession(prevSession.GetId(), session)
		}
	} else {
		slog.Debug("[MQTT] Starting new session from running", "client-id", clientId)
		session.SetConnected(true)
		if id, err := persistence.SessionRepository.PersistSession(session); err != nil {
			slog.Error("[MQTT] error persisting clean session", "client-id", clientId, "err", err)
		} else {
			session.ApplySessionId(id)
		}
	}
}
