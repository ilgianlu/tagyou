package mqtt

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ilgianlu/tagyou/event"
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
	"nhooyr.io/websocket"
)

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
			return
		}
		defer c.Close(websocket.StatusInternalError, "the sky is falling")

		events := make(chan *packet.Packet)
		go event.RangeEvents(connections, events)

		ctx, cancel := context.WithTimeout(r.Context(), time.Second*10)
		defer cancel()

		msgType, msg, err := c.Read(ctx)

		if err != nil {
			log.Err(err).Msg("error reading message")
			return
		}

		log.Debug().Msg(fmt.Sprintf("received type %s : %s", msgType.String(), string(msg)))

		c.Close(websocket.StatusNormalClosure, "")
	}

}
