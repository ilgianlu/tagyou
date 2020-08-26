package packet

import (
	"testing"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
)

func TestConnectSuccess(t *testing.T) {
	e := make(chan Packet)
	s := model.Session{}
	p := Packet{
		header:          16,
		remainingLength: 20,
		remainingBytes:  []byte{0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 0, 0, 7, 99, 108, 105, 101, 110, 116, 88},
		Session:         &s,
	}
	go connectSuccessEvent(t, e)
	p.connectReq()
}

func connectSuccessEvent(t *testing.T, e chan Packet) {
	packet := <-e
	if packet.Event != EVENT_CONNECT {
		t.Errorf("expected event type %d, found %d", EVENT_CONNECT, packet.Event)
	}
	if packet.Session.ProtocolVersion != conf.MQTT_V5 {
		t.Errorf("expected protocol version %d, found %d", conf.MQTT_V5, packet.Session.ProtocolVersion)
	}
}

func TestConnectSuccessWithProperties(t *testing.T) {
	e := make(chan Packet)
	s := model.Session{}
	p := Packet{
		header:          16,
		remainingLength: 25,
		remainingBytes:  []byte{0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 5, 17, 0, 0, 0, 60, 0, 7, 99, 108, 105, 101, 110, 116, 88},
		Session:         &s,
	}
	go connectSuccessEventWithProperties(t, e)
	p.connectReq()
}

func connectSuccessEventWithProperties(t *testing.T, e chan Packet) {
	packet := <-e
	if packet.Event != EVENT_CONNECT {
		t.Errorf("expected event type %d, found %d", EVENT_CONNECT, packet.Event)
	}
	if packet.Session.ProtocolVersion != conf.MQTT_V5 {
		t.Errorf("expected protocol version %d, found %d", conf.MQTT_V5, packet.Session.ProtocolVersion)
	}
	if len(packet.properties) != 1 {
		t.Errorf("expected 1 property, found %d", len(packet.properties))
	}
	if packet.SessionExpiryInterval() != 60 {
		t.Errorf("expected sessione expiry interval 60 sec, found %d", packet.SessionExpiryInterval())
	}
}
