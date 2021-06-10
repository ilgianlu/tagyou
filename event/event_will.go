package event

import (
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/ec5"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	kgo "github.com/segmentio/kafka-go"
	"gorm.io/gorm"
)

func sendWill(db *gorm.DB, kwriter *kgo.Writer, p *packet.Packet, outQueue chan<- *out.OutData) {
	if p.Session.WillTopic != "" {
		willPacket := packet.Publish(
			p.Session.ProtocolVersion,
			p.Session.WillQoS(),
			p.Session.WillRetain(),
			p.Session.WillTopic,
			packet.NewPacketIdentifier(),
			p.Session.WillMessage,
		)
		if conf.KAFKA_ON {
			ec5.Publish(kwriter, p.Session.WillTopic, &willPacket)
		}
		sendForward(db, p.Session.WillTopic, &willPacket, outQueue)
	}
}
