package event

import (
	"context"
	"database/sql"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ilgianlu/tagyou/conf"
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
	mockConn := TagyouConnMock{}
	session := model.RunningSession{
		Connected: false,
		Conn:      mockConn,
		Router:    routers.NewSimple(),
	}
	p := packet.Publish(4, 0, false, "", 0, []byte{})

	managePacket(&session, &p)

	if session.Connected != false {
		t.Errorf("expecting not connected session false, received true")
	}
}

func TestConnectWhenClientIsAlreadyConnected(t *testing.T) {
	mockConn := TagyouConnMock{}
	session := model.RunningSession{
		Connected: true,
		ClientId:  "client-x",
		Conn:      mockConn,
		Router:    routers.NewSimple(),
	}
	p := packet.Connect()

	managePacket(&session, &p)

	if session.Connected != false {
		t.Errorf("expecting not connected session false, received true")
	}
}

func TestSuccessfullConnect(t *testing.T) {
	mockConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)},
	}
	session := model.RunningSession{
		Connected: false,
		Conn:      mockConn,
		Router:    routers.NewSimple(),
	}

	buf := []byte{16, 25, 0, 4, 77, 81, 84, 84, 5, 2, 0, 30, 5, 17, 0, 0, 0, 60, 0, 7, 99, 108, 105, 101, 110, 116, 88}
	p, err := packet.PacketParse(&session, buf)

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

	router := routers.NewSimple()

	receivedMsgs = []Msg{}
	mockConn := TagyouConnMock{
		remoteAddr: &net.IPAddr{IP: net.IPv4(127, 0, 0, 1)},
	}
	session := model.RunningSession{
		Connected: false,
		Conn:      mockConn,
		Router:    router,
	}

	buf := []byte{16, 64, 0, 4, 77, 81, 84, 84, 4, 198, 0, 5, 0, 15, 109, 113, 116, 116, 106, 115, 95, 97, 97, 50, 51, 99, 56, 49, 53, 0, 5, 97, 47, 98, 47, 99, 0, 15, 119, 105, 108, 108, 32, 109, 101, 115, 115, 97, 103, 101, 46, 46, 46, 0, 4, 117, 115, 101, 114, 0, 5, 112, 108, 117, 116, 111}
	p, _ := packet.PacketParse(&session, buf)

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
