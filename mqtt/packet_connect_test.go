package mqtt

import (
	"testing"

	"github.com/ilgianlu/tagyou/model"
)

func TestConnect(t *testing.T) {
	e := make(chan Event)
	s := model.Session{}
	p := Packet{
		header:          16,
		remainingLength: 20,
		remainingBytes:  []byte{0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 0, 0, 7, 99, 108, 105, 101, 110, 116, 88},
	}
	go connectSuccessEvent(t, e)
	connectReq(p, e, &s)
}

func connectSuccessEvent(t *testing.T, e chan Event) {
	event := <-e
	if event.err != 0 {
		t.Errorf("did not expect any error, found %d", event.err)
	}
	if event.eventType != EVENT_CONNECT {
		t.Errorf("expected event type %d, found %d", EVENT_CONNECT, event.eventType)
	}
	if event.session.ProtocolVersion != MQTT_V5 {
		t.Errorf("expected protocol version %d, found %d", MQTT_V5, event.session.ProtocolVersion)
	}
}
