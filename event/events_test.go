package event

import (
	"net"
	"testing"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
	"github.com/ilgianlu/tagyou/sqlrepository"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Msg []byte

type MsgQueue []Msg

var receivedMsgs []Msg

type TagyouConnMock struct {
	remoteAddr net.Addr
}

func (mock TagyouConnMock) Write(b []byte) (n int, err error) {
	receivedMsgs = append(receivedMsgs, b)
	return len(b), nil
}

func (mock TagyouConnMock) Close() error {
	return nil
}

func (mock TagyouConnMock) RemoteAddr() net.Addr {
	return mock.remoteAddr
}

func TestPublishWhenClientIsNotConnected(t *testing.T) {
	router := routers.NewSimple()
	mockConn := TagyouConnMock{}
	session := model.RunningSession{
		Connected: false,
		Conn:      mockConn,
	}
	p := packet.Packet{Event: packet.EVENT_PUBLISH}

	manageEvent(router, &session, &p)

	if session.Connected != false {
		t.Errorf("expecting not connected session false, received true")
	}
}

func TestConnectWhenClientIsAlreadyConnected(t *testing.T) {
	router := routers.NewSimple()
	mockConn := TagyouConnMock{}
	session := model.RunningSession{
		Connected: true,
		ClientId:  "client-x",
		Conn:      mockConn,
	}
	p := packet.Packet{Event: packet.EVENT_CONNECT}

	manageEvent(router, &session, &p)

	if session.Connected != false {
		t.Errorf("expecting not connected session false, received true")
	}
}

func TestSuccessfullConnect(t *testing.T) {
	router := routers.NewSimple()
	mockConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)},
	}
	session := model.RunningSession{
		Connected: false,
		Conn:      mockConn,
	}

	buf := []byte{16, 25, 0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 5, 17, 0, 0, 0, 60, 0, 7, 99, 108, 105, 101, 110, 116, 88}
	p, err := packet.PacketParse(&session, buf)

	if err != nil {
		t.Errorf("did not expect parse error %s", err)
	}

	manageEvent(router, &session, &p)

	if session.Connected != true {
		t.Errorf("expecting client connected, received false")
	}

	if session.GetClientId() != "clientX" {
		t.Errorf("expecting clientId clientX, received %s", session.GetClientId())
	}
}

func TestSuccessfullReconnect(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	persistence := persistence.SqlPersistence{}
	persistence.InnerInit(db, false, false)

	storedSession := sqlrepository.Session{
		ID:              5,
		LastSeen:        time.Now().Unix() - 1000,
		LastConnect:     time.Now().Unix() - 2000,
		ExpiryInterval:  3600,
		ClientId:        "mqttjs_aa23c815",
		Connected:       false,
		ProtocolVersion: conf.MQTT_V3_11,
	}
	db.Save(&storedSession)

	router := routers.NewSimple()

	receivedMsgs = []Msg{}
	mockConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)},
	}
	session := model.RunningSession{
		Connected: false,
		Conn:      mockConn,
	}

	buf := []byte{16, 64, 0, 4, 77, 81, 84, 84, 4, 198, 0, 5, 0, 15, 109, 113, 116, 116, 106, 115, 95, 97, 97, 50, 51, 99, 56, 49, 53, 0, 5, 97, 47, 98, 47, 99, 0, 15, 119, 105, 108, 108, 32, 109, 101, 115, 115, 97, 103, 101, 46, 46, 46, 0, 4, 117, 115, 101, 114, 0, 5, 112, 108, 117, 116, 111}
	p, _ := packet.PacketParse(&session, buf)

	manageEvent(router, &session, &p)

	if session.Connected != true {
		t.Errorf("expecting client connected, received false")
	}

	if session.GetClientId() != "mqttjs_aa23c815" {
		t.Errorf("expecting clientId mqttjs_aa23c815, received %s", session.GetClientId())
	}

	if !session.CleanStart() {
		t.Errorf("expected mqttjs_aa23c815 clean start, received false")
	}

	if !router.DestinationExists("mqttjs_aa23c815") {
		t.Errorf("expected mqttjs_aa23c815 to exist in router, received false")
	}

	if len(receivedMsgs) != 1 {
		t.Errorf("expected 1 msg received in mqttjs_aa23c815, received %d", len(receivedMsgs))
	}

	if receivedMsgs[0][0] != 32 {
		t.Errorf("expected connack, received %d", receivedMsgs[0][0])
	}
}

func TestPublish(t *testing.T) {
	conf.Loader()
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	persistence := persistence.SqlPersistence{}
	persistence.InnerInit(db, false, false)

	db.Exec("DELETE FROM sessions")
	db.Exec("DELETE FROM subscriptions")

	router := routers.NewSimple()

	receivedMsgs = []Msg{}

	// publisher
	publisherConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(10, 0, 0, 1)},
	}
	router.AddDestination("publisher", publisherConn)
	publisherSession := model.RunningSession{
		Connected: true,
		ClientId:  "publisher",
		Conn:      publisherConn,
	}
	pubStoredSession := sqlrepository.Session{
		ID:              5,
		LastSeen:        time.Now().Unix() - 2000,
		LastConnect:     time.Now().Unix() - 2000,
		ExpiryInterval:  3600,
		ClientId:        "publisher",
		Connected:       true,
		ProtocolVersion: conf.MQTT_V3_11,
	}
	db.Save(&pubStoredSession)

	// subscriber
	subscriberConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(10, 0, 0, 2)},
	}
	router.AddDestination("subscriber", subscriberConn)
	subscriberSession := model.RunningSession{
		Connected: true,
		ClientId:  "subscriber",
		Conn:      subscriberConn,
	}
	subStoredSession := sqlrepository.Session{
		ID:              6,
		LastSeen:        time.Now().Unix() - 2000,
		LastConnect:     time.Now().Unix() - 2000,
		ExpiryInterval:  3600,
		ClientId:        "subscriber",
		Connected:       true,
		ProtocolVersion: conf.MQTT_V3_11,
	}
	db.Save(&subStoredSession)

	topic := "a-topic-to-subscribe"
	sub := sqlrepository.Subscription{
		Topic:     topic,
		ClientId:  "subscriber",
		SessionID: 6,
	}
	db.Save(&sub)

	p := packet.Publish(4, 0, false, topic, 0, []byte("pippo"))
	p.Event = packet.EVENT_PUBLISH
	p.Topic = topic

	manageEvent(router, &publisherSession, &p)

	if !router.DestinationExists("publisher") {
		t.Errorf("expected publisher to exist in router, received false")
	}

	if !router.DestinationExists("subscriber") {
		t.Errorf("expected subscriber to exist in router, received false")
	}

	if len(receivedMsgs) != 1 {
		t.Errorf("expected 1 msg received by subscriber, received %d", len(receivedMsgs))
	}

	receivedPacket, _ := packet.PacketParse(&subscriberSession, receivedMsgs[0])
	if receivedPacket.PacketType() != packet.PACKET_TYPE_PUBLISH {
		t.Errorf("expected %d (publish) msg received, received %d", packet.PACKET_TYPE_PUBLISH, receivedPacket.PacketType())
	}
}

func TestSubscribe(t *testing.T) {
	conf.Loader()
	db, err := gorm.Open(sqlite.Open("test.db3"), &gorm.Config{})
	if err != nil {
		log.Fatal().Err(err).Msg("[API] failed to connect database")
	}

	persistence := persistence.SqlPersistence{}
	persistence.InnerInit(db, false, false)

	db.Exec("DELETE FROM sessions")
	db.Exec("DELETE FROM subscriptions")

	router := routers.NewSimple()

	receivedMsgs = []Msg{}

	// subscriber
	subscriberConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(10, 0, 0, 2)},
	}
	router.AddDestination("client", subscriberConn)
	subscriberSession := model.RunningSession{
		Connected: true,
		ClientId:  "client",
		Conn:      subscriberConn,
	}
	subStoredSession := sqlrepository.Session{
		ID:              6,
		LastSeen:        time.Now().Unix() - 2000,
		LastConnect:     time.Now().Unix() - 2000,
		ExpiryInterval:  3600,
		ClientId:        "client",
		Connected:       true,
		ProtocolVersion: conf.MQTT_V3_11,
	}
	db.Save(&subStoredSession)

	// (0x82) subscription of client 'client' to topic '/topic/#'
	buf := []byte{0x82, 0x0d, 0x33, 0x41, 0x00, 0x08, 0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x2f, 0x23, 0x00}

	p, _ := packet.PacketParse(&subscriberSession, buf)

	manageEvent(router, &subscriberSession, &p)

	if len(receivedMsgs) != 1 {
		t.Errorf("expected 1 msg received by subscriber, received %d", len(receivedMsgs))
	}

	receivedPacket, _ := packet.PacketParse(&subscriberSession, receivedMsgs[0])
	if receivedPacket.PacketType() != packet.PACKET_TYPE_SUBACK {
		t.Errorf("expected %d (publish) msg received, received %d", packet.PACKET_TYPE_PUBLISH, receivedPacket.PacketType())
	}

	subscription := sqlrepository.Subscription{}
	if err := db.Where("client_id = ?", "client").First(&subscription).Error; err != nil {
		t.Errorf("unexpected error during find of sub %s", err)
	}

	if subscription.Topic != "/topic/#" {
		t.Errorf("expecting subscription of %s, found %s", "/topic/#", subscription.Topic)
	}

	if subscription.ClientId != "client" {
		t.Errorf("expecting clientId subscription of %s, found %s", "client", subscription.ClientId)
	}
}
