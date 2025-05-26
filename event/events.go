package event

import (
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func saveRetain(session *model.RunningSession, p *packet.Packet) {
	var r model.Retain
	r.ClientID = session.ClientId
	r.Topic = p.Topic
	r.ApplicationMessage = p.ApplicationMessage()
	r.CreatedAt = time.Now().Unix()
	persistence.RetainRepository.Delete(r)
	if len(r.ApplicationMessage) > 0 {
		persistence.RetainRepository.Create(r)
	}
}
