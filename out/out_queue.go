package out

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"gorm.io/gorm"
)

func RangeOutQueue(connections *model.Connections, db *gorm.DB, outQueue <-chan OutData) {
	for o := range outQueue {
		simpleSend(connections, db, o.ClientId, o.Packet)
	}
}

func simpleSend(connections *model.Connections, db *gorm.DB, clientId string, p []byte) {
	conn, exists := connections.Exists(clientId)
	if exists {
		if conn == nil {
			log.Error().Msgf("cannot write to %s net.Conn, c is nil (removing)", clientId)
			connections.Remove(clientId)
			return
		}
		// packetBytes := p.ToByteSlice()
		_, err := conn.Write(p)
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
