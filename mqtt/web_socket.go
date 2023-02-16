package mqtt

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/event"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

type WebsocketConnection struct {
	conn websocket.Conn
}

func (c WebsocketConnection) Write(b []byte) (n int, e error) {
	err := c.conn.Write(context.Background(), websocket.MessageBinary, b)
	return len(b), err
}

func (c WebsocketConnection) Close() error {
	return c.conn.Close(websocket.StatusNormalClosure, "")
}

func (c WebsocketConnection) RemoteAddr() net.Addr {
	return &net.IPAddr{}
}

func StartWebSocket(port string, connections *model.Connections) {
	r := httprouter.New()
	r.GET("/ws", AcceptWebsocket(connections))

	log.Info().Msgf("[WS] websocket listening on %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal().Err(err).Msg("[WS] websocket listener broken")
	}
}

func AcceptWebsocket(connections *model.Connections) func(http.ResponseWriter, *http.Request, httprouter.Params) {

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			Subprotocols: []string{"mqtt"},
		})
		if err != nil {
			log.Err(err).Msg("error accepting websocket connection")
			c.Close(websocket.StatusInternalError, "the sky is falling")
			return
		}

		events := make(chan *packet.Packet)
		go event.RangeEvents(connections, events)

		session := model.RunningSession{
			KeepAlive:      conf.DEFAULT_KEEPALIVE,
			ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
			Conn:           WebsocketConnection{conn: *c},
			LastConnect:    time.Now().Unix(),
		}

		bytesFromWs := make(chan []byte)
		defer close(bytesFromWs)

		go readFromWs(&session, r, c, bytesFromWs)
		handleMqtt(&session, bytesFromWs, events)
	}

}

func readFromWs(session *model.RunningSession, r *http.Request, c *websocket.Conn, bytesFromWs chan<- []byte) {
	for {
		// log.Debug().Msg(fmt.Sprintf("keep alive set %d secs", session.GetKeepAlive()*2))
		// ctx, cancel := context.WithTimeout(r.Context(), time.Duration(session.GetKeepAlive()*2)*time.Second)

		msgType, msg, err := c.Read(context.Background())
		if err != nil {
			log.Err(err).Msg("error reading message")
			// cancel()
			return
		}

		log.Debug().Msg(fmt.Sprintf("received type %s : %s", msgType.String(), string(msg)))
		bytesFromWs <- msg
		// cancel()
	}
}

func handleMqtt(session *model.RunningSession, bytesFromWs <-chan []byte, events chan<- *packet.Packet) {
	buf := []byte{}
	for bytesRead := range bytesFromWs {
		log.Debug().Msg(fmt.Sprintf("received message : %s", string(bytesRead)))

		if len(bytesRead) == 0 {
			continue
		}

		buf = append(buf, bytesRead...)

		pb, err := packet.ReadFromByteSlice(buf)
		if err != nil {
			log.Debug().Msgf("error during ReadFromByteSlice : %s", err.Error())
			continue
		}

		buf = buf[len(buf):]

		p, err := packetParse(session, pb)
		if err != nil {
			log.Debug().Msgf("error during packet parse : %s", err.Error())
			return
		}

		events <- &p

	}
}
