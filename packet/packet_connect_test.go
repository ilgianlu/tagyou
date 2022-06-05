package packet

import (
	"testing"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/rs/zerolog"
)

func TestConnectSuccess(t *testing.T) {
	s := model.RunningSession{}
	p := Packet{
		header:          16,
		remainingLength: 20,
		remainingBytes:  []byte{0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 0, 0, 7, 99, 108, 105, 101, 110, 116, 88},
		Session:         &s,
	}
	p.connectReq()
	if p.Event != EVENT_CONNECT {
		t.Errorf("expected event type %d, found %d", EVENT_CONNECT, p.Event)
	}
	if p.Session.ProtocolVersion != conf.MQTT_V5 {
		t.Errorf("expected protocol version %d, found %d", conf.MQTT_V5, p.Session.ProtocolVersion)
	}
}

func TestConnectSuccessWithProperties(t *testing.T) {
	s := model.RunningSession{}
	p := Packet{
		header:          16,
		remainingLength: 25,
		remainingBytes:  []byte{0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 5, 17, 0, 0, 0, 60, 0, 7, 99, 108, 105, 101, 110, 116, 88},
		Session:         &s,
	}
	p.connectReq()
	if p.Event != EVENT_CONNECT {
		t.Errorf("expected event type %d, found %d", EVENT_CONNECT, p.Event)
	}
	if p.Session.ProtocolVersion != conf.MQTT_V5 {
		t.Errorf("expected protocol version %d, found %d", conf.MQTT_V5, p.Session.ProtocolVersion)
	}
	if len(p.properties) != 1 {
		t.Errorf("expected 1 property, found %d", len(p.properties))
	}
	if p.SessionExpiryInterval() != 60 {
		t.Errorf("expected sessione expiry interval 60 sec, found %d", p.SessionExpiryInterval())
	}
}

func BenchmarkStartConnect(b *testing.B) {
	buffer := []byte{16, 59, 0, 4, 77, 81, 84, 84, 4, 6, 0, 60, 0, 15, 109, 113, 116, 116, 106, 115, 95, 53, 51, 48, 48, 102, 100, 54, 51, 0, 8, 108, 97, 115, 116, 119, 105, 108, 108, 0, 20, 97, 32, 118, 101, 114, 121, 32, 115, 104, 111, 114, 116, 32, 109, 101, 115, 115, 97, 103, 101}
	for n := 0; n < b.N; n++ {
		Start(buffer)
	}
}

func BenchmarkStartSubscribe(b *testing.B) {
	buffer := []byte{130, 13, 149, 223, 0, 8, 112, 114, 101, 115, 101, 110, 99, 101, 0}
	for n := 0; n < b.N; n++ {
		Start(buffer)
	}
}

func BenchmarkStartPublish(b *testing.B) {
	buffer := []byte{48, 20, 0, 8, 112, 114, 101, 115, 101, 110, 99, 101, 72, 101, 108, 108, 111, 32, 109, 113, 116, 116}
	for n := 0; n < b.N; n++ {
		Start(buffer)
	}
}

func BenchmarkParseConnect(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	buffer := []byte{16, 59, 0, 4, 77, 81, 84, 84, 4, 6, 0, 60, 0, 15, 109, 113, 116, 116, 106, 115, 95, 53, 51, 48, 48, 102, 100, 54, 51, 0, 8, 108, 97, 115, 116, 119, 105, 108, 108, 0, 20, 97, 32, 118, 101, 114, 121, 32, 115, 104, 111, 114, 116, 32, 109, 101, 115, 115, 97, 103, 101}
	p, _ := Start(buffer)

	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		LastConnect:    time.Now().Unix(),
	}
	p.Session = &session

	for n := 0; n < b.N; n++ {
		p.Parse()
	}
}

func BenchmarkParseSubscribe(b *testing.B) {
	buffer := []byte{130, 13, 149, 223, 0, 8, 112, 114, 101, 115, 101, 110, 99, 101, 0}
	p, _ := Start(buffer)

	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		LastConnect:    time.Now().Unix(),
	}
	p.Session = &session

	for n := 0; n < b.N; n++ {
		p.Parse()
	}
}

func BenchmarkParsePublish(b *testing.B) {
	buffer := []byte{48, 20, 0, 8, 112, 114, 101, 115, 101, 110, 99, 101, 72, 101, 108, 108, 111, 32, 109, 113, 116, 116}
	p, _ := Start(buffer)

	session := model.RunningSession{
		KeepAlive:      conf.DEFAULT_KEEPALIVE,
		ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
		LastConnect:    time.Now().Unix(),
	}
	p.Session = &session

	for n := 0; n < b.N; n++ {
		p.Parse()
	}
}
