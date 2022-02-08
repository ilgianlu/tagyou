package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/gorm"
)

func onConnect(db *gorm.DB, connections *model.Connections, p *packet.Packet, outQueue chan<- out.OutData) {
	clientId := p.Session.GetClientId()
	if conf.FORBID_ANONYMOUS_LOGIN && !p.Session.FromLocalhost() {
		if !doAuth(db, p.Session) {
			return
		}
	}
	taken := checkConnectionTakeOver(p, connections, outQueue)
	if taken {
		log.Debug().Msgf("%s reconnecting", clientId)
	}
	connections.Add(clientId, p.Session.GetConn())

	startSession(db, p.Session)

	connack := packet.Connack(false, packet.CONNECT_OK, p.Session.GetProtocolVersion())
	sendSimple(clientId, &connack, outQueue)
}

func doAuth(db *gorm.DB, session *model.RunningSession) bool {
	session.Mu.RLock()
	clientId := session.ClientId
	username := session.Username
	password := session.Password
	session.Mu.RUnlock()
	ok, pubAcl, subAcl := model.CheckAuth(db, clientId, username, password)
	if !ok {
		log.Debug().Msg("wrong connect credentials")
		return false
	}
	session.ApplyAcl(pubAcl, subAcl)
	log.Debug().Msgf("auth ok, imported acls %s, %s", pubAcl, subAcl)
	return true
}

func checkConnectionTakeOver(p *packet.Packet, connections *model.Connections, outQueue chan<- out.OutData) bool {
	p.Session.Mu.RLock()
	clientId := p.Session.ClientId
	protocolVersion := p.Session.ProtocolVersion
	p.Session.Mu.RUnlock()
	if _, ok := connections.Exists(clientId); !ok {
		return false
	}

	pkt := packet.Connack(false, packet.SESSION_TAKEN_OVER, protocolVersion)
	sendSimple(clientId, &pkt, outQueue)

	err := connections.Close(clientId)
	if err != nil {
		log.Debug().Msgf("%s : error taking over another connection; %s", clientId, err)
	}
	connections.Remove(clientId)
	return true
}

func startSession(db *gorm.DB, session *model.RunningSession) {
	clientId := session.GetClientId()
	if prevSession, ok := model.SessionExists(db, clientId); ok {
		if session.CleanStart() || prevSession.Expired() || session.GetProtocolVersion() != prevSession.ProtocolVersion {
			log.Debug().Msgf("%s Cleaning previous session: Invalid or to clean", clientId)
			if err := model.CleanSession(db, clientId); err != nil {
				log.Err(err).Msgf("%s : error removing previous session", clientId)
			}
			if id, err := model.PersistSession(db, session, true); err != nil {
				log.Err(err).Msgf("%s : error persisting clean session", clientId)
			} else {
				session.ApplySessionId(id)
			}
		} else {
			log.Debug().Msgf("%s Updating previous session from running", clientId)
			session.ApplySessionId(prevSession.ID)
			prevSession.UpdateFromRunning(session)
			db.Save(&prevSession)
		}
	} else {
		log.Debug().Msgf("%s Starting new session from running", clientId)
		if id, err := model.PersistSession(db, session, true); err != nil {
			log.Err(err).Msgf("%s : error persisting clean session", clientId)
		} else {
			session.ApplySessionId(id)
		}
	}
}
