package out

import (
	"log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/jinzhu/gorm"
)

func RangeOutQueue(connections model.Connections, db *gorm.DB, outQueue <-chan *OutData) {
	for o := range outQueue {
		simpleSend(connections, db, o.ClientId, o.Packet)
	}
}

func simpleSend(connections model.Connections, db *gorm.DB, clientId string, p packet.Packet) {
	if c, ok := connections[clientId]; ok {
		n, err := c.Write(p.ToByteSlice())
		if err != nil {
			log.Println("cannot write to", clientId, ":", err)
		} else {
			log.Println("published", n, "bytes to", clientId)
		}
	} else {
		log.Println(clientId, "is not connected")
	}
}
