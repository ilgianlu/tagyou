package event

import (
	"time"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
)

func saveRetain(clientId string, topic string, applicationMessage []byte) {
	var r model.Retain
	r.ClientID = clientId
	r.Topic = topic
	r.ApplicationMessage = applicationMessage
	r.CreatedAt = time.Now().Unix()
	persistence.RetainRepository.Delete(r)
	if len(r.ApplicationMessage) > 0 {
		persistence.RetainRepository.Create(r)
	}
}
