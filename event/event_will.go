package event

import (
	"github.com/ilgianlu/tagyou/nowherecloud"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"gorm.io/gorm"
)

func sendWill(db *gorm.DB, p *packet.Packet, ncMessages chan<- nowherecloud.NcMessage, outQueue chan<- out.OutData) {
	if p.Session.WillTopic != "" {
		willPacket := packet.Publish(
			p.Session.ProtocolVersion,
			p.Session.WillQoS(),
			p.Session.WillRetain(),
			p.Session.WillTopic,
			packet.NewPacketIdentifier(),
			p.Session.WillMessage,
		)
		if nowherecloud.KAFKA_ON {
			// nowherecloud.Publish(kwriter, p.Session.WillTopic, &willPacket)
			ncMessages <- nowherecloud.NcMessage{Topic: p.Session.WillTopic, P: &willPacket}
		}
		sendForward(db, p.Session.WillTopic, &willPacket, outQueue)
	}
}
