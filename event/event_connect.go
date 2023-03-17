package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sender"
)

func onConnect(sender sender.Sender, p *packet.Packet) {
	clientId := p.Session.GetClientId()
	if conf.FORBID_ANONYMOUS_LOGIN && !p.Session.FromLocalhost() {
		if !doAuth(p.Session) {
			return
		}
	}
	taken := checkConnectionTakeOver(p, sender)
	if taken {
		log.Debug().Msgf("[MQTT] (%s) reconnecting", clientId)
	}
	sender.AddDestination(clientId, p.Session.GetConn())

	startSession(p.Session)

	connack := packet.Connack(false, packet.CONNECT_OK, p.Session.GetProtocolVersion())
	sender.Send(clientId, connack.ToByteSlice())
}

func doAuth(session *model.RunningSession) bool {
	session.Mu.RLock()
	clientId := session.ClientId
	username := session.Username
	password := session.Password
	session.Mu.RUnlock()
	ok, pubAcl, subAcl := checkAuth(clientId, username, password)
	if !ok {
		log.Debug().Msg("[MQTT] wrong connect credentials")
		return false
	}
	session.ApplyAcl(pubAcl, subAcl)
	log.Debug().Msgf("[MQTT] auth ok, imported acls %s, %s", pubAcl, subAcl)
	return true
}

func checkAuth(clientId string, username string, password string) (bool, string, string) {
	auth, err := persistence.AuthRepository.GetByClientIdUsername(clientId, username)
	if err != nil {
		log.Debug().Msg("[MQTT] could not find user")
		return false, "", ""
	}

	mAuth := model.Auth(auth)
	if err := mAuth.CheckPassword(password); err != nil {
		log.Debug().Msg("[MQTT] wrong password")
		return false, "", ""
	}

	return true, auth.PublishAcl, auth.SubscribeAcl
}

func checkConnectionTakeOver(p *packet.Packet, sender sender.Sender) bool {
	p.Session.Mu.RLock()
	clientId := p.Session.ClientId
	protocolVersion := p.Session.ProtocolVersion
	p.Session.Mu.RUnlock()
	if sender.DestinationExists(clientId) {
		return false
	}

	pkt := packet.Connack(false, packet.SESSION_TAKEN_OVER, protocolVersion)
	sender.Send(clientId, pkt.ToByteSlice())

	sender.RemoveDestination(clientId)
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
