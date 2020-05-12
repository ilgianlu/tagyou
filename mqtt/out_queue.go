package mqtt

import (
	"log"
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/jinzhu/gorm"
)

func rangeOutQueue(connections Connections, db *gorm.DB, outQueue <-chan OutData) {
	for o := range outQueue {
		simpleSend(connections, db, o.clientId, o.packet)
	}
}

func simpleSend(connections Connections, db *gorm.DB, clientId string, packet Packet) {
	if c, ok := connections[clientId]; ok {
		n, err := c.Write(packet.toByteSlice())
		if err != nil {
			log.Println("cannot write to", clientId, ":", err)
		} else {
			log.Println("published", n, "bytes to", clientId)
		}
		if packet.QoS() == 1 || packet.QoS() == 2 {
			saveRetry(db, clientId, packet)
		}
	} else {
		log.Println(clientId, "is not connected")
	}
}

func saveRetry(db *gorm.DB, clientId string, packet Packet) {
	var r model.Retry
	r.ClientId = clientId
	r.PacketIdentifier = packet.packetIdentifier
	r.Qos = packet.QoS()
	r.Dup = packet.Dup()
	if r.Qos == 1 {
		r.AckStatus = model.WAIT_FOR_PUB_ACK
	} else {
		r.AckStatus = model.WAIT_FOR_PUB_REC
	}
	r.ApplicationMessage = packet.ApplicationMessage()
	r.CreatedAt = time.Now()
	db.Create(&r)
}
