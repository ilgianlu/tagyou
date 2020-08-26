package websocket

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/time/rate"
	ws "nhooyr.io/websocket"
)

func (wc WebsocketController) UpgradeConnection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, err := ws.Accept(w, r, &ws.AcceptOptions{})
	if err != nil {
		log.Printf("[WS] err %s", err)
		return
	}
	defer c.Close(ws.StatusInternalError, "[WS] the sky is falling")

	l := rate.NewLimiter(rate.Every(time.Second*1), 10)
	for {
		err = echo(r.Context(), c, l)
		if ws.CloseStatus(err) == ws.StatusNormalClosure {
			return
		}
		if err != nil {
			log.Printf("failed to echo with %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

func echo(ctx context.Context, c *ws.Conn, l *rate.Limiter) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*100)
	defer cancel()

	err := l.Wait(ctx)
	if err != nil {
		return err
	}

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	if err != nil {
		return fmt.Errorf("[WS] failed to io.Copy: %w", err)
	}
	log.Println("[WS] echoed")

	err = w.Close()
	return err
}
