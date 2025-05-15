package mqtt

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/coder/websocket"
	"github.com/ilgianlu/tagyou/api/controllers/middlewares"
	"github.com/ilgianlu/tagyou/conf"
	"github.com/ilgianlu/tagyou/event"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
)

func StartWebSocket(port string, router model.Router) {
	r := http.NewServeMux()
	r.HandleFunc("GET /ws", AcceptWebsocket(router))
	r.HandleFunc("POST /messages", middlewares.Authenticated(PostMessage(router)))

	slog.Info("[WS] websocket listening on", "tcp-port", port)
	if err := http.ListenAndServe(port, r); err != nil {
		slog.Error("[WS] websocket listener broken", "err", err)
		panic(1)
	}
}

func AcceptWebsocket(router model.Router) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			Subprotocols: []string{"mqtt"},
		})
		if err != nil {
			slog.Error("[WS] error accepting websocket connection", "err", err)
			c.Close(websocket.StatusInternalError, "the sky is falling")
			return
		}

		session := model.RunningSession{
			KeepAlive:      conf.DEFAULT_KEEPALIVE,
			ExpiryInterval: int64(conf.SESSION_MAX_DURATION_SECONDS),
			Conn:           model.WebsocketConnection{Conn: c},
			LastConnect:    time.Now().Unix(),
		}

		events := make(chan *packet.Packet)
		go event.RangePackets(router, &session, events)

		bytesFromWs := make(chan []byte)
		defer close(bytesFromWs)

		go readFromWs(c, bytesFromWs)
		handleMqtt(&session, bytesFromWs, events)
	}
}

func PostMessage(router model.Router) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		mess := model.Message{}
		if err := json.NewDecoder(r.Body).Decode(&mess); err != nil {
			slog.Error("[WS] error decoding json input", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		session := model.RunningSession{ClientId: "PostMessage"}

		msg := packet.Publish(4, mess.Qos, mess.Retained, mess.Topic, 0, payloadFromPayloadType(mess.Payload))
		msg.Topic = mess.Topic
		event.OnPublish(router, &session, &msg)

		if res, err := json.Marshal("message published"); err != nil {
			slog.Error("[WS] error marshaling response message", "err", err)
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
			slog.Info("[WS] Wrote json result", "num-bytes", numBytes)
		}
	}
}

func payloadFromPayloadType(payload string) []byte {
	return []byte(payload)
}

func readFromWs(c *websocket.Conn, bytesFromWs chan<- []byte) {
	for {
		msgType, msg, err := c.Read(context.Background())
		if err != nil {
			slog.Error("[WS] error reading message", "err", err)
			return
		}

		slog.Debug("[WS] received message", "type", msgType.String(), "msg", string(msg))
		bytesFromWs <- msg
	}
}

func handleMqtt(session *model.RunningSession, bytesFromWs <-chan []byte, events chan<- *packet.Packet) {
	buf := []byte{}
	for bytesRead := range bytesFromWs {
		slog.Debug("[WS] received message", "msg-bytes", string(bytesRead))

		if len(bytesRead) == 0 {
			continue
		}

		buf = append(buf, bytesRead...)

		pb, err := packet.ReadFromByteSlice(buf)
		if err != nil {
			slog.Debug("[WS] error during ReadFromByteSlice", "err", err.Error())
			continue
		}

		buf = buf[len(buf):]

		p, err := packet.PacketParse(session, pb)
		if err != nil {
			slog.Debug("[WS] error during packet parse", "err", err.Error())
			return
		}

		events <- &p
	}
}
