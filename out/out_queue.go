package out

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/gorm"
)

func RangeOutQueue(connections model.Connections, db *gorm.DB, outQueue <-chan *OutData) {
	for o := range outQueue {
		simpleSend(connections, db, o.ClientId, o.Packet)
	}
}

func simpleSend(connections model.Connections, db *gorm.DB, clientId string, p packet.Packet) {
	if c, ok := connections[clientId]; ok {
		_, err := c.Write(p.ToByteSlice())
		if err != nil {
			log.Debug().Err(err).Msgf("cannot write to %s", clientId)
		}
		// else {
		// 	log.Println("published", n, "bytes to", clientId)
		// }
	} else {
		log.Debug().Msgf("%s is not connected", clientId)
	}
}
