package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/jinzhu/gorm"
)

func onConnect(db *gorm.DB, connections Connections, p packet.Packet, outQueue chan<- OutData) {
	if conf.FORBID_ANONYMOUS_LOGIN && !p.Session.FromLocalhost() {
		ok, pubAcl, subAcl := model.CheckAuth(db, p.Session.ClientId, p.Session.Username, p.Session.Password)
		if !ok {
			log.Println("wrong connect credentials")
			return
		} else {
			p.Session.PublishAcl = pubAcl
			p.Session.SubscribeAcl = subAcl
			log.Printf("auth ok, imported acls %s, %s\n", pubAcl, subAcl)
		}
	}
	taken := checkConnectionTakeOver(p, connections, outQueue)
	if taken {
		log.Printf("%s reconnecting", p.Session.ClientId)
	}
	connections.Add(p.Session.ClientId, p.Session.Conn)

	startSession(db, p.Session)

	sendSimple(p.Session.ClientId, packet.Connack(false, packet.CONNECT_OK, p.Session.ProtocolVersion), outQueue)
}

func checkConnectionTakeOver(p packet.Packet, connections Connections, outQueue chan<- OutData) bool {
	if _, ok := connections.Exists(p.Session.ClientId); !ok {
		return false
	}

	pkt := packet.Connack(false, packet.SESSION_TAKEN_OVER, p.Session.ProtocolVersion)
	sendSimple(p.Session.ClientId, pkt, outQueue)

	err := connections.Close(p.Session.ClientId)
	if err != nil {
		log.Printf("%s : error taking over another connection; %s", p.Session.ClientId, err)
	}
	connections.Remove(p.Session.ClientId)
	return true
}

func startSession(db *gorm.DB, session *model.Session) {
	if prevSession, ok := model.SessionExists(db, session.ClientId); ok {
		if session.CleanStart() || prevSession.Expired() || session.ProtocolVersion != prevSession.ProtocolVersion {
			model.CleanSession(db, session.ClientId)
		}
		prevSession.MergeSession(*session)
		session = &prevSession
		db.Save(&session)
	} else {
		db.Create(&session)
	}
}
