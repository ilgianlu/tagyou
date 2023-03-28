package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ilgianlu/tagyou/api/controllers/middlewares"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/event"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/routers"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

func StartWebSocket(port string, router routers.Router) {
	r := httprouter.New()
	r.GET("/ws", AcceptWebsocket(router))
	r.POST("/messages", middlewares.Authenticated(PostMessage(router)))

	log.Info().Msgf("[WS] websocket listening on %s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatal().Err(err).Msg("[WS] websocket listener broken")
	}
}

func AcceptWebsocket(router routers.Router) func(http.ResponseWriter, *http.Request, httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			Subprotocols: []string{"mqtt"},
		})
		if err != nil {
			log.Err(err).Msg("error accepting websocket connection")
			c.Close(websocket.StatusInternalError, "the sky is falling")
			return
		}

		session := model.RunningSession{
			KeepAlive:      conf.DEFAULT_KEEPALIVE,
			ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
			Conn:           model.WebsocketConnection{Conn: *c},
			LastConnect:    time.Now().Unix(),
		}

		events := make(chan *packet.Packet)
		go event.RangeEvents(router, &session, events)

		bytesFromWs := make(chan []byte)
		defer close(bytesFromWs)

		go readFromWs(r, c, bytesFromWs)
		handleMqtt(&session, bytesFromWs, events)
	}
}

func PostMessage(router routers.Router) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		mess := model.Message{}
		if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
			log.Printf("error decoding json input: %s\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		session := model.RunningSession{ClientId: "PostMessage"}

		msg := packet.Publish(4, mess.Qos, mess.Retained, mess.Topic, 0, payloadFromPayloadType(mess.Payload, mess.PayloadType))
		msg.Topic = mess.Topic
		event.OnPublish(router, &session, &msg)

		if res, err := json.Marshal("message published"); err != nil {
			log.Printf("error marshaling response message: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			numBytes, err := w.Write(res)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			log.Printf("Wrote %d bytes json result\n", numBytes)
		}
	}
}

func payloadFromPayloadType(payload string, payloadType byte) []byte {
	return []byte(payload)
}

func readFromWs(r *http.Request, c *websocket.Conn, bytesFromWs chan<- []byte) {
	for {
		msgType, msg, err := c.Read(context.Background())
		if err != nil {
			log.Err(err).Msg("error reading message")
			return
		}

		log.Debug().Msg(fmt.Sprintf("received type %s : %s", msgType.String(), string(msg)))
		bytesFromWs <- msg
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

		p, err := packet.PacketParse(session, pb)
		if err != nil {
			log.Debug().Msgf("error during packet parse : %s", err.Error())
			return
		}

		events <- &p
	}
}
