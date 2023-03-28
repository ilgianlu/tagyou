package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/password"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func onConnect(router routers.Router, session *model.RunningSession, p *packet.Packet) {
	clientId := session.GetClientId()
	if conf.FORBID_ANONYMOUS_LOGIN && !session.FromLocalhost() {
		if !doAuth(session) {
			return
		}
	}
	taken := checkConnectionTakeOver(session, router)
	if taken {
		log.Debug().Msgf("[MQTT] (%s) reconnecting", clientId)
	}
	router.AddDestination(clientId, session.GetConn())

	startSession(session)

	connack := packet.Connack(false, packet.CONNECT_OK, session.GetProtocolVersion())
	router.Send(clientId, connack.ToByteSlice())
}

func doAuth(session *model.RunningSession) bool {
	session.Mu.RLock()
	clientId := session.ClientId
	username := session.Username
	sessionPassword := session.Password
	session.Mu.RUnlock()
	ok, pubAcl, subAcl := checkAuth(clientId, username, sessionPassword)
	if !ok {
		log.Debug().Msg("[MQTT] wrong connect credentials")
		return false
	}
	session.ApplyAcl(pubAcl, subAcl)
	log.Debug().Msgf("[MQTT] auth ok, imported acls %s, %s", pubAcl, subAcl)
	return true
}

func checkAuth(clientId string, username string, sessionPassword string) (bool, string, string) {
	client, err := persistence.ClientRepository.GetByClientIdUsername(clientId, username)
	if err != nil {
		log.Debug().Msg("[MQTT] could not find user")
		return false, "", ""
	}

	if err := password.CheckPassword(client.Password, []byte(sessionPassword)); err != nil {
		log.Debug().Msg("[MQTT] wrong password")
		return false, "", ""
	}

	return true, client.PublishAcl, client.SubscribeAcl
}

func checkConnectionTakeOver(session *model.RunningSession, router routers.Router) bool {
	session.Mu.RLock()
	clientId := session.ClientId
	protocolVersion := session.ProtocolVersion
	session.Mu.RUnlock()
	if !router.DestinationExists(clientId) {
		return false
	}

	pkt := packet.Connack(false, packet.SESSION_TAKEN_OVER, protocolVersion)
	router.Send(clientId, pkt.ToByteSlice())

	router.RemoveDestination(clientId)
	return true
}

func startSession(session *model.RunningSession) {
	clientId := session.GetClientId()
	if prevSession, ok := persistence.SessionRepository.SessionExists(clientId); ok {
		if session.CleanStart() || prevSession.Expired() || session.GetProtocolVersion() != prevSession.GetProtocolVersion() {
			log.Debug().Msgf("[MQTT] check session (%t) (%t) (%d != %d)", session.CleanStart(), prevSession.Expired(), session.GetProtocolVersion(), prevSession.GetProtocolVersion())
			log.Debug().Msgf("[MQTT] (%s) Cleaning previous session: Invalid or to clean", clientId)
			if err := persistence.SessionRepository.CleanSession(clientId); err != nil {
				log.Err(err).Msgf("[MQTT] (%s) error removing previous session", clientId)
			}
			session.SetConnected(true)
			if id, err := persistence.SessionRepository.PersistSession(session); err != nil {
				log.Err(err).Msgf("[MQTT] (%s) error persisting clean session", clientId)
			} else {
				session.ApplySessionId(id)
			}
		} else {
			log.Debug().Msgf("%s Updating previous session from running", clientId)
			session.ApplySessionId(prevSession.GetId())
			persistence.SessionRepository.Save(&prevSession)
		}
	} else {
		log.Debug().Msgf("[MQTT] (%s) Starting new session from running", clientId)
		session.SetConnected(true)
		if id, err := persistence.SessionRepository.PersistSession(session); err != nil {
			log.Err(err).Msgf("[MQTT] (%s) error persisting clean session", clientId)
		} else {
			session.ApplySessionId(id)
		}
	}
}
