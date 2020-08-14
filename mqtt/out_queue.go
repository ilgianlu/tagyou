package mqtt

import (
	"log"

	"github.com/jinzhu/gorm"
)

func rangeOutQueue(connections Connections, db *gorm.DB, outQueue <-chan OutData) {
	for o := range outQueue {
		simpleSend(connections, db, o.clientId, o.packet)
	}
}

func simpleSend(connections Connections, db *gorm.DB, clientId string, packet Packet) {
	if c, ok := connections[clientId]; ok {
		b := packet.ToByteSlice()
		log.Println(b)
		n, err := c.Write(b)
		if err != nil {
			log.Println("cannot write to", clientId, ":", err)
		} else {
			log.Println("published", n, "bytes to", clientId)
		}
	} else {
		log.Println(clientId, "is not connected")
	}
}
