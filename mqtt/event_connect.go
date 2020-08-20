package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func onConnect(db *gorm.DB, connections Connections, p Packet, outQueue chan<- OutData) {
	if conf.FORBID_ANONYMOUS_LOGIN && !p.session.FromLocalhost() {
		ok, pubAcl, subAcl := model.CheckAuth(db, p.session.ClientId, p.session.Username, p.session.Password)
		if !ok {
			log.Println("wrong connect credentials")
			return
		} else {
			p.session.PublishAcl = pubAcl
			p.session.SubscribeAcl = subAcl
			log.Printf("auth ok, imported acls %s, %s\n", pubAcl, subAcl)
		}
	}
	taken := checkConnectionTakeOver(p, connections, outQueue)
	if taken {
		log.Printf("%s reconnecting", p.session.ClientId)
	}
	connections.Add(p.session.ClientId, p.session.Conn)

	startSession(db, p.session)

	sendSimple(p.session.ClientId, Connack(false, CONNECT_OK, p.session.ProtocolVersion), outQueue)
}

func checkConnectionTakeOver(p Packet, connections Connections, outQueue chan<- OutData) bool {
	if _, ok := connections.Exists(p.session.ClientId); !ok {
		return false
	}

	pkt := Connack(false, SESSION_TAKEN_OVER, p.session.ProtocolVersion)
	sendSimple(p.session.ClientId, pkt, outQueue)

	err := connections.Close(p.session.ClientId)
	if err != nil {
		log.Printf("%s : error taking over another connection; %s", p.session.ClientId, err)
	}
	connections.Remove(p.session.ClientId)
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
