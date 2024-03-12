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
	"github.com/ilgianlu/tagyou/sqlc"
	"github.com/ilgianlu/tagyou/sqlc/dbaccess"
)

func TestSubscribe(t *testing.T) {
	conf.Loader()
	os.Setenv("DEBUG", "1")
	os.Remove("test.db3")

	dbConn, err := sql.Open("sqlite3", "test.db3")
	if err != nil {
		t.Errorf("[API] failed to connect database")
	}

	dbConn.ExecContext(context.Background(), "PRAGMA foreign_keys = ON;")
	dbConn.ExecContext(context.Background(), sqlc.DBSchema)

	db := dbaccess.New(dbConn)

	persistence := persistence.SqlPersistence{}
	persistence.InnerInit(db, false, false, []byte(""))

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
