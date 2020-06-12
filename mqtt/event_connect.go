package mqtt

import (
	"log"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func onConnect(db *gorm.DB, connections Connections, e Event, outQueue chan<- OutData) {
	if conf.DISALLOW_ANONYMOUS_LOGIN && !model.CheckAuth(db, e.clientId, e.session.Username, e.session.Password) {
		log.Println("wrong connect credentials")
		return
	}

	if c, ok := connections[e.clientId]; ok {
		log.Println("session taken over")
		p := Connack(false, SESSION_TAKEN_OVER, e.session.ProtocolVersion)
		sendSimple(e.clientId, p, outQueue)
		closeClient(c)
		removeClient(e.clientId, connections)
	}
	connections[e.clientId] = e.session.Conn
	sendSimple(e.clientId, Connack(false, CONNECT_OK, e.session.ProtocolVersion), outQueue)

	startSession(db, e.session)
}

func startSession(db *gorm.DB, session *model.Session) {
	if db.Where("client_id = ?", session.ClientId).First(&session).RecordNotFound() {
		db.Create(&session)
	} else {
		if session.CleanStart() {
			model.CleanSession(db, session.ClientId)
			db.Create(&session)
		} else {
			db.Save(&session)
		}
	}
}
