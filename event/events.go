package event

import (
	"time"

	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sender"
)

func RangeEvents(sender sender.Sender, events <-chan *packet.Packet) {
	for p := range events {
		clientId := p.Session.GetClientId()
		switch p.Event {
		case packet.EVENT_CONNECT:
			log.Debug().Msgf("//!! EVENT type %d client connect %s", p.Event, clientId)
			onConnect(sender, p)
		case packet.EVENT_SUBSCRIBED:
			log.Debug().Msgf("//!! EVENT type %d client subscribed %s", p.Event, clientId)
			onSubscribe(sender, p)
		case packet.EVENT_UNSUBSCRIBED:
			log.Debug().Msgf("//!! EVENT type %d client unsubscribed %s", p.Event, clientId)
			onUnsubscribe(sender, p)
		case packet.EVENT_PUBLISH:
			log.Debug().Msgf("//!! EVENT type %d client published to %s %s QoS %d", p.Event, p.Topic, clientId, p.QoS())
			OnPublish(sender, p)
		case packet.EVENT_PUBACKED:
			log.Debug().Msgf("//!! EVENT type %d client acked message %d %s", p.Event, p.PacketIdentifier(), clientId)
			clientPuback(p)
		case packet.EVENT_PUBRECED:
			log.Debug().Msgf("//!! EVENT type %d pub received message %d %s", p.Event, p.PacketIdentifier(), clientId)
			clientPubrec(sender, p)
		case packet.EVENT_PUBRELED:
			log.Debug().Msgf("//!! EVENT type %d pub releases message %d %s", p.Event, p.PacketIdentifier(), clientId)
			clientPubrel(sender, p)
		case packet.EVENT_PUBCOMPED:
			log.Debug().Msgf("//!! EVENT type %d pub complete message %d %s", p.Event, p.PacketIdentifier(), clientId)
			clientPubcomp(p)
		case packet.EVENT_PING:
			log.Debug().Msgf("//!! EVENT type %d client ping %s", p.Event, clientId)
			onPing(sender, p)
		case packet.EVENT_DISCONNECT:
			log.Debug().Msgf("//!! EVENT type %d client disconnect %s", p.Event, clientId)
			clientDisconnect(sender, p, clientId)
		case packet.EVENT_WILL_SEND:
			log.Debug().Msgf("//!! EVENT type %d sending will message %s", p.Event, clientId)
			sendWill(sender, p)
		case packet.EVENT_PACKET_ERR:
			log.Debug().Msgf("//!! EVENT type %d packet error %s", p.Event, clientId)
			clientDisconnect(sender, p, clientId)
		}
	}
}

func onPing(sender sender.Sender, p *packet.Packet) {
	toSend := packet.PingResp()
	sender.Send(p.Session.GetClientId(), toSend.ToByteSlice())
}

func clientDisconnect(sender sender.Sender, p *packet.Packet, clientId string) {
	if sender.DestinationExists(clientId) {
		needDisconnection := needDisconnection(p)
		if !needDisconnection {
			return
		}
		sender.RemoveDestination(clientId)
		persistence.SessionRepository.DisconnectSession(clientId)
	}
}

func saveRetain(p *packet.Packet) {
	var r model.Retain
	r.Topic = p.Topic
	r.ApplicationMessage = p.ApplicationMessage()
	r.CreatedAt = time.Now().Unix()
	persistence.RetainRepository.Delete(r)
	if len(r.ApplicationMessage) > 0 {
		persistence.RetainRepository.Create(r)
	}
}

func needDisconnection(p *packet.Packet) bool {
	if session, ok := persistence.SessionRepository.SessionExists(p.Session.ClientId); ok {
		log.Debug().Msgf("[MQTT] (%s) Persisted session LastConnect %d running session %d", p.Session.ClientId, session.LastConnect, p.Session.LastConnect)
		if session.LastConnect > p.Session.LastConnect {
			// session persisted is newer then running memory session... device reconnected!
			// no need to send will
			log.Debug().Msgf("[MQTT] (%s) avoid disconnect! (device reconnected)", p.Session.ClientId)
			return false
		}
	}
	return true
}
