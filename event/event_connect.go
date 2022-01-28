package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/gorm"
)

func onConnect(db *gorm.DB, connections model.Connections, p *packet.Packet, outQueue chan<- *out.OutData) {
	if conf.FORBID_ANONYMOUS_LOGIN && !p.Session.FromLocalhost() {
		ok, pubAcl, subAcl := model.CheckAuth(db, p.Session.ClientId, p.Session.Username, p.Session.Password)
		if !ok {
			log.Debug().Msg("wrong connect credentials")
			return
		} else {
			p.Session.PublishAcl = pubAcl
			p.Session.SubscribeAcl = subAcl
			log.Debug().Msgf("auth ok, imported acls %s, %s", pubAcl, subAcl)
		}
	}
	taken := checkConnectionTakeOver(p, connections, outQueue)
	if taken {
		log.Debug().Msgf("%s reconnecting", p.Session.ClientId)
	}
	connections.Add(p.Session.ClientId, p.Session.Conn)

	startSession(db, p.Session)

	connack := packet.Connack(false, packet.CONNECT_OK, p.Session.ProtocolVersion)
	sendSimple(p.Session.ClientId, &connack, outQueue)
}

func checkConnectionTakeOver(p *packet.Packet, connections model.Connections, outQueue chan<- *out.OutData) bool {
	if _, ok := connections.Exists(p.Session.ClientId); !ok {
		return false
	}

	pkt := packet.Connack(false, packet.SESSION_TAKEN_OVER, p.Session.ProtocolVersion)
	sendSimple(p.Session.ClientId, &pkt, outQueue)

	err := connections.Close(p.Session.ClientId)
	if err != nil {
		log.Debug().Msgf("%s : error taking over another connection; %s", p.Session.ClientId, err)
	}
	connections.Remove(p.Session.ClientId)
	return true
}

func startSession(db *gorm.DB, session *model.RunningSession) {
	if prevSession, ok := model.SessionExists(db, session.ClientId); ok {
		if session.CleanStart() || prevSession.Expired() || session.ProtocolVersion != prevSession.ProtocolVersion {
			log.Debug().Msgf("%s Cleaning previous session: Invalid or to clean", session.ClientId)
			if err := model.CleanSession(db, session.ClientId); err != nil {
				log.Err(err).Msgf("%s : error removing previous session", session.ClientId)
			}
			if id, err := model.PersistSession(db, session, true); err != nil {
				log.Err(err).Msgf("%s : error persisting clean session", session.ClientId)
			} else {
				session.SessionID = id
			}
		} else {
			log.Debug().Msgf("%s Updating previous session from running", session.ClientId)
			session.SessionID = prevSession.ID
			prevSession.UpdateFromRunning(*session)
			db.Save(&prevSession)
		}
	} else {
		log.Debug().Msgf("%s Starting new session from running", session.ClientId)
		if id, err := model.PersistSession(db, session, true); err != nil {
			log.Err(err).Msgf("%s : error persisting clean session", session.ClientId)
		} else {
			session.SessionID = id
		}
	}
}
