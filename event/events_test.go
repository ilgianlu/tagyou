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

type TagyouConnMock struct {
	remoteAddr net.Addr
}

func (mock TagyouConnMock) Write(b []byte) (n int, err error) {
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
	persistence.InnerInit(db)

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
		t.Errorf("expected clientX clean start, received false")
	}
}
