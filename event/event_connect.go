package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func onConnect(connections *model.Connections, p *packet.Packet) {
	clientId := p.Session.GetClientId()
	if conf.FORBID_ANONYMOUS_LOGIN && !p.Session.FromLocalhost() {
		if !doAuth(p.Session) {
			return
		}
	}
	taken := checkConnectionTakeOver(p, connections)
	if taken {
		log.Debug().Msgf("[MQTT] (%s) reconnecting", clientId)
	}
	connections.Add(clientId, p.Session.GetConn())

	startSession(p.Session)

	connack := packet.Connack(false, packet.CONNECT_OK, p.Session.GetProtocolVersion())
	SimpleSend(connections, clientId, connack.ToByteSlice())
}

func doAuth(session *model.RunningSession) bool {
	session.Mu.RLock()
	clientId := session.ClientId
	username := session.Username
	password := session.Password
	session.Mu.RUnlock()
	ok, pubAcl, subAcl := CheckAuth(clientId, username, password)
	if !ok {
		log.Debug().Msg("[MQTT] wrong connect credentials")
		return false
	}
	session.ApplyAcl(pubAcl, subAcl)
	log.Debug().Msgf("[MQTT] auth ok, imported acls %s, %s", pubAcl, subAcl)
	return true
}

func CheckAuth(clientId string, username string, password string) (bool, string, string) {
	auth, err := persistence.AuthRepository.GetByClientIdUsername(clientId, username)
	if err != nil {
		return false, "", ""
	}

	mAuth := model.Auth(auth)
	if err := mAuth.CheckPassword(password); err != nil {
		return false, "", ""
	}

	return true, auth.PublishAcl, auth.SubscribeAcl
}

func checkConnectionTakeOver(p *packet.Packet, connections *model.Connections) bool {
	p.Session.Mu.RLock()
	clientId := p.Session.ClientId
	protocolVersion := p.Session.ProtocolVersion
	p.Session.Mu.RUnlock()
	if _, ok := connections.Exists(clientId); !ok {
		return false
	}

	pkt := packet.Connack(false, packet.SESSION_TAKEN_OVER, protocolVersion)
	SimpleSend(connections, clientId, pkt.ToByteSlice())

	err := connections.Close(clientId)
	if err != nil {
		log.Debug().Msgf("[MQTT] (%s) error taking over another connection : %s", clientId, err)
	}
	connections.Remove(clientId)
	return true
}

func startSession(session *model.RunningSession) {
	clientId := session.GetClientId()
	if prevSession, ok := persistence.SessionRepository.SessionExists(clientId); ok {
		if session.CleanStart() || prevSession.Expired() || session.GetProtocolVersion() != prevSession.ProtocolVersion {
			log.Debug().Msgf("[MQTT] (%s) Cleaning previous session: Invalid or to clean", clientId)
			if err := persistence.SessionRepository.CleanSession(clientId); err != nil {
				log.Err(err).Msgf("[MQTT] (%s) error removing previous session", clientId)
			}
			if id, err := persistence.SessionRepository.PersistSession(session, true); err != nil {
				log.Err(err).Msgf("[MQTT] (%s) error persisting clean session", clientId)
			} else {
				session.ApplySessionId(id)
			}
		} else {
			log.Debug().Msgf("%s Updating previous session from running", clientId)
			session.ApplySessionId(prevSession.ID)
			prevSession.UpdateFromRunning(session)
			persistence.SessionRepository.Save(&prevSession)
		}
	} else {
		log.Debug().Msgf("[MQTT] (%s) Starting new session from running", clientId)
		if id, err := persistence.SessionRepository.PersistSession(session, true); err != nil {
			log.Err(err).Msgf("[MQTT] (%s) error persisting clean session", clientId)
		} else {
			session.ApplySessionId(id)
		}
	}
}
