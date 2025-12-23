package mqtt

import (
	"bufio"
	"bytes"
	"context"
	"database/sql"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/engine"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
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
	connections := model.SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	mockConn := TagyouConnMock{}
	session := model.RunningSession{
		Connected: false,
		Conn:      mockConn,
		Router:    routers.NewSimple(&connections),
		Engine:    engine.NewEngine(),
	}
	p := packet.Publish(4, 0, false, "", 0, []byte{})

	managePacket(&session, &p)

	if session.Connected != false {
		t.Errorf("expecting not connected session false, received true")
	}
}

func TestConnectWhenClientIsAlreadyConnected(t *testing.T) {
	connections := model.SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	mockConn := TagyouConnMock{}
	session := model.RunningSession{
		Connected: true,
		ClientId:  "client-x",
		Conn:      mockConn,
		Router:    routers.NewSimple(&connections),
		Engine:    engine.NewEngine(),
	}
	p := packet.Connect()

	managePacket(&session, &p)

	if session.Connected != false {
		t.Errorf("expecting not connected session false, received true")
	}
}

func TestSuccessfullConnect(t *testing.T) {
	ps := persistence.SqlPersistence{DbFile: "test.db3", InitDatabase: true}
	_, err := ps.Init(false, false, []byte{})
	if err != nil {
		t.Errorf("did not expect any error opening test.db3")
	}

	connections := model.SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	mockConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)},
	}
	session := model.RunningSession{
		Connected: false,
		Conn:      mockConn,
		Router:    routers.NewSimple(&connections),
		Engine:    engine.NewEngine(),
	}

	buf := bytes.NewReader([]byte{16, 25, 0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 5, 17, 0, 0, 0, 60, 0, 7, 99, 108, 105, 101, 110, 116, 88})
	p := packet.Packet{}
	err = p.Parse(bufio.NewReader(buf), &session)
	if err != nil {
		t.Errorf("did not expect parse error %s", err)
	}

	managePacket(&session, &p)

	if session.Connected != true {
		t.Errorf("expecting client connected, received false")
	}

	if session.GetClientId() != "clientX" {
		t.Errorf("expecting clientId clientX, received %s", session.GetClientId())
	}
}

func TestSuccessfullReconnect(t *testing.T) {
	os.Setenv("DEBUG", "1")
	ps := persistence.SqlPersistence{DbFile: "test.db3", InitDatabase: true}
	db, err := ps.Init(false, false, []byte{})
	if err != nil {
		t.Errorf("did not expect any error opening test.db3")
	}

	db.CreateSession(context.Background(), dbaccess.CreateSessionParams{
		LastSeen:        sql.NullInt64{Int64: time.Now().Unix() - 1000, Valid: true},
		LastConnect:     sql.NullInt64{Int64: time.Now().Unix() - 2000, Valid: true},
		ExpiryInterval:  sql.NullInt64{Int64: 3600, Valid: true},
		ClientID:        sql.NullString{String: "mqttjs_aa23c815", Valid: true},
		Connected:       sql.NullInt64{Int64: 0, Valid: true},
		ProtocolVersion: sql.NullInt64{Int64: conf.MQTT_V3_11, Valid: true},
	})

	connections := model.SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	router := routers.NewSimple(&connections)

	receivedMsgs = []Msg{}
	mockConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)},
	}
	session := model.RunningSession{
		Connected: false,
		Conn:      mockConn,
		Router:    router,
		Engine:    engine.NewEngine(),
	}

	buf := bytes.NewReader([]byte{16, 64, 0, 4, 77, 81, 84, 84, 4, 198, 0, 5, 0, 15, 109, 113, 116, 116, 106, 115, 95, 97, 97, 50, 51, 99, 56, 49, 53, 0, 5, 97, 47, 98, 47, 99, 0, 15, 119, 105, 108, 108, 32, 109, 101, 115, 115, 97, 103, 101, 46, 46, 46, 0, 4, 117, 115, 101, 114, 0, 5, 112, 108, 117, 116, 111})
	p := packet.Packet{}
	err = p.Parse(bufio.NewReader(buf), &session)

	managePacket(&session, &p)

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
	os.Setenv("DEBUG", "1")

	ps := persistence.SqlPersistence{DbFile: "test.db3", InitDatabase: true}
	db, err := ps.Init(false, false, []byte{})
	if err != nil {
		t.Errorf("did not expect any error opening test.db3")
	}

	connections := model.SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	router := routers.NewSimple(&connections)

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
		Router:    router,
		Engine:    engine.NewEngine(),
	}

	db.CreateSession(context.Background(), dbaccess.CreateSessionParams{
		LastSeen:        sql.NullInt64{Int64: time.Now().Unix() - 2000, Valid: true},
		LastConnect:     sql.NullInt64{Int64: time.Now().Unix() - 2000, Valid: true},
		ExpiryInterval:  sql.NullInt64{Int64: 3600, Valid: true},
		ClientID:        sql.NullString{String: "publisher", Valid: true},
		Connected:       sql.NullInt64{Int64: 1, Valid: true},
		ProtocolVersion: sql.NullInt64{Int64: conf.MQTT_V3_11, Valid: true},
	})

	// subscriber
	subscriberConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(10, 0, 0, 2)},
	}
	router.AddDestination("subscriber", subscriberConn)
	subscriberSession := model.RunningSession{
		Connected: true,
		ClientId:  "subscriber",
		Conn:      subscriberConn,
		Engine:    engine.NewEngine(),
	}

	subSession, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{
		LastSeen:        sql.NullInt64{Int64: time.Now().Unix() - 2000, Valid: true},
		LastConnect:     sql.NullInt64{Int64: time.Now().Unix() - 2000, Valid: true},
		ExpiryInterval:  sql.NullInt64{Int64: 3600, Valid: true},
		ClientID:        sql.NullString{String: "subscriber", Valid: true},
		Connected:       sql.NullInt64{Int64: 1, Valid: true},
		ProtocolVersion: sql.NullInt64{Int64: conf.MQTT_V3_11, Valid: true},
	})

	topic := "a-topic-to-subscribe"
	db.CreateSubscription(context.Background(), dbaccess.CreateSubscriptionParams{
		ShareName: sql.NullString{String: "", Valid: true},
		Shared:    sql.NullInt64{Int64: 0, Valid: true},
		Topic:     sql.NullString{String: topic, Valid: true},
		ClientID:  sql.NullString{String: "subscriber", Valid: true},
		SessionID: sql.NullInt64{Int64: subSession.ID, Valid: true},
	})

	p := packet.Publish(4, 0, false, topic, 0, []byte("pippo"))

	managePacket(&publisherSession, &p)

	if !router.DestinationExists("publisher") {
		t.Errorf("expected publisher to exist in router, received false")
	}

	if !router.DestinationExists("subscriber") {
		t.Errorf("expected subscriber to exist in router, received false")
	}

	if len(receivedMsgs) != 1 {
		t.Errorf("expected 1 msg received by subscriber, received %d", len(receivedMsgs))
	}

	// receivedPacket, _ := packet.PacketParse(&subscriberSession, receivedMsgs[0])
	buf := bytes.NewReader(receivedMsgs[0])
	p = packet.Packet{}
	err = p.Parse(bufio.NewReader(buf), &subscriberSession)
	if p.PacketType() != packet.PACKET_TYPE_PUBLISH {
		t.Errorf("expected %d (publish) msg received, received %d", packet.PACKET_TYPE_PUBLISH, p.PacketType())
	}
}

func TestSubscribe(t *testing.T) {
	conf.Loader()
	os.Setenv("DEBUG", "1")

	ps := persistence.SqlPersistence{DbFile: "test.db3", InitDatabase: true}
	db, err := ps.Init(false, false, []byte{})
	if err != nil {
		t.Errorf("did not expect any error opening test.db3")
	}

	connections := model.SimpleConnections{
		Conns: make(map[string]model.TagyouConn, conf.ROUTER_STARTING_CAPACITY),
	}
	router := routers.NewSimple(&connections)

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
		Router:    router,
		Engine:    engine.NewEngine(),
	}
	sess, _ := db.CreateSession(context.Background(), dbaccess.CreateSessionParams{
		LastSeen:        sql.NullInt64{Int64: time.Now().Unix() - 2000, Valid: true},
		LastConnect:     sql.NullInt64{Int64: time.Now().Unix() - 2000, Valid: true},
		ExpiryInterval:  sql.NullInt64{Int64: 3600, Valid: true},
		ClientID:        sql.NullString{String: "client", Valid: true},
		Connected:       sql.NullInt64{Int64: 1, Valid: true},
		ProtocolVersion: sql.NullInt64{Int64: conf.MQTT_V3_11, Valid: true},
	})
	subscriberSession.SessionID = sess.ID

	// (0x82) subscription of client 'client' to topic '/topic/#'
	buf := bytes.NewReader([]byte{0x82, 0x0d, 0x33, 0x41, 0x00, 0x08, 0x2f, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x2f, 0x23, 0x00})

	p := packet.Packet{}
	err = p.Parse(bufio.NewReader(buf), &subscriberSession)

	managePacket(&subscriberSession, &p)

	if len(receivedMsgs) != 1 {
		t.Errorf("expected 1 msg received by subscriber, received %d", len(receivedMsgs))
	}

	// receivedPacket, _ := tParse(&subscriberSession, receivedMsgs[0])
	buf = bytes.NewReader(receivedMsgs[0])
	p = packet.Packet{}
	err = p.Parse(bufio.NewReader(buf), &subscriberSession)
	if p.PacketType() != packet.PACKET_TYPE_SUBACK {
		t.Errorf("expected %d (publish) msg received, received %d", packet.PACKET_TYPE_PUBLISH, p.PacketType())
	}

	subscriptions, err := db.GetSubscriptionsBySessionId(context.Background(), sql.NullInt64{Int64: sess.ID, Valid: true})
	if err != nil {
		t.Errorf("unexpected error during find of sub %s", err)
	}

	if len(subscriptions) == 0 {
		t.Errorf("expected 1 subscription, received %d", len(subscriptions))
		return
	}

	if subscriptions[0].Topic.String != "/topic/#" {
		t.Errorf("expecting subscription of %s, found %s", "/topic/#", subscriptions[0].Topic.String)
	}

	if subscriptions[0].ClientID.String != "client" {
		t.Errorf("expecting clientId subscription of %s, found %s", "client", subscriptions[0].ClientID.String)
	}
}
