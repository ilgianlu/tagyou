package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func onConnect(db *gorm.DB, connections Connections, e Event, outQueue chan<- OutData) {
	if conf.FORBID_ANONYMOUS_LOGIN && !model.CheckAuth(db, e.clientId, e.session.Username, e.session.Password) {
		log.Println("wrong connect credentials")
		return
	}
	taken := checkConnectionTakeOver(e, connections, outQueue)
	if taken {
		log.Printf("%s reconnecting", e.clientId)
	}
	connections.Add(e.clientId, e.session.Conn)

	startSession(db, e.session)

	sendSimple(e.clientId, Connack(false, CONNECT_OK, e.session.ProtocolVersion), outQueue)
}

func checkConnectionTakeOver(e Event, connections Connections, outQueue chan<- OutData) bool {
	if _, ok := connections.Exists(e.clientId); !ok {
		return false
	}

	p := Connack(false, SESSION_TAKEN_OVER, e.session.ProtocolVersion)
	sendSimple(e.clientId, p, outQueue)

	err := connections.Close(e.clientId)
	if err != nil {
		log.Printf("%s : error taking over another connection; %s", e.clientId, err)
	}
	connections.Remove(e.clientId)
	return true
}

func startSession(db *gorm.DB, session *model.Session) {
	if prevSession, ok := model.SessionExists(db, session.ClientId); ok {
		if session.CleanStart() || prevSession.Expired() {
			model.CleanSession(db, session.ClientId)
		}
		prevSession.MergeSession(*session)
		session = &prevSession
		db.Save(&session)
	} else {
		db.Create(&session)
	}
}
