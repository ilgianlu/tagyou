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

func TestPublish(t *testing.T) {
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
