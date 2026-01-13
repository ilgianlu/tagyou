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
	clientID := session.GetClientId()
	isAuth := doAuth(session)
	if conf.FORBID_ANONYMOUS_LOGIN && !isAuth {
		return
	}
	taken := checkConnectionTakeOver(session)
	if taken {
		slog.Debug("[MQTT] client reconnecting", "client-id", clientID)
	}
	session.Router.AddDestination(clientID, session.GetConn())

	startSession(session)

	connack := packet.Connack(false, packet.CONNECT_OK, session.GetProtocolVersion())
	bs, err := connack.ToByteSlice()
	if err != nil {
		return
	}
	session.Router.Send(clientID, bs)
}

func doAuth(session *model.RunningSession) bool {
	clientID := session.ClientId
	username := session.Username
	if username == "" {
		slog.Debug("[MQTT] client did not pass any username")
		return false
	}
	sessionPassword := session.Password
	if sessionPassword == "" {
		slog.Debug("[MQTT] client did not pass any password")
		return false
	}
	ok, pubACL, subACL := checkAuth(clientID, username, sessionPassword)
	if !ok {
		slog.Debug("[MQTT] wrong connect credentials")
		return false
	}
	session.ApplyAcl(pubACL, subACL)
	slog.Debug("[MQTT] auth ok, imported acls", "pub-acl", pubACL, "sub-acl", subACL)
	return true
}

func checkAuth(clientID string, username string, sessionPassword string) (bool, string, string) {
	client, err := persistence.ClientRepository.GetByClientIdUsername(clientID, username)
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
	clientID := session.ClientId
	protocolVersion := session.ProtocolVersion
	if !session.Router.DestinationExists(clientID) {
		return false
	}

	pkt := packet.Connack(false, packet.SESSION_TAKEN_OVER, protocolVersion)
	bs, err := pkt.ToByteSlice()
	if err != nil {
		return false
	}
	session.Router.Send(clientID, bs)

	session.Router.RemoveDestination(clientID)
	return true
}

func startSession(session *model.RunningSession) {
	clientID := session.GetClientId()
	session.Router = routers.ByClientId(clientID, session.Router.GetConns())
	if prevSession, ok := persistence.SessionRepository.SessionExists(clientID); ok {
		slog.Debug("[MQTT] check existing session", "last-seen", prevSession.GetLastSeen(), "clean-start", session.CleanStart(), "expired", prevSession.Expired(), "new-protocol-version", session.GetProtocolVersion(), "prev-protocol-version", prevSession.GetProtocolVersion())
		if session.CleanStart() || prevSession.Expired() || session.GetProtocolVersion() != prevSession.GetProtocolVersion() {
			slog.Debug("[MQTT] Cleaning previous session: Invalid or to clean", "client-id", clientID)
			if err := persistence.SessionRepository.CleanSession(clientID); err != nil {
				slog.Error("[MQTT] error removing previous session", "client-id", clientID, "err", err)
			}
			session.SetConnected(true)
			if id, err := persistence.SessionRepository.PersistSession(session); err != nil {
				slog.Error("[MQTT] error persisting clean session", "client-id", clientID, "err", err)
			} else {
				session.ApplySessionId(id)
			}
		} else {
			slog.Debug("Updating previous session from running", "client-id", clientID)
			session.ApplySessionId(prevSession.GetId())
			session.SetConnected(true)
			_, _ = persistence.SessionRepository.UpdateSession(prevSession.GetId(), session)
		}
	} else {
		slog.Debug("[MQTT] Starting new session from running", "client-id", clientID)
		session.SetConnected(true)
		if id, err := persistence.SessionRepository.PersistSession(session); err != nil {
			slog.Error("[MQTT] error persisting clean session", "client-id", clientID, "err", err)
		} else {
			session.ApplySessionId(id)
		}
	}
}
