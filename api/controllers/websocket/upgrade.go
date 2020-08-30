package websocket

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/ilgianlu/tagyou/packet"
	"github.com/julienschmidt/httprouter"
	ws "nhooyr.io/websocket"
)

func (wc WebsocketController) UpgradeConnection(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	c, err := ws.Accept(w, r, &ws.AcceptOptions{})
	if err != nil {
		log.Printf("[WS] err %s", err)
		return
	}
	defer c.Close(ws.StatusInternalError, "[WS] the sky is falling")

	buf := bytes.NewBuffer([]byte{})
	go readInto(r, c, buf)

	scanner := bufio.NewScanner(buf)
	packetSplit := func(b []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(b) == 0 && atEOF {
			// socket down - closed
			return 0, b, bufio.ErrFinalToken
		}
		pb, err := packet.ReadFromByteSlice(b)
		if err != nil {
			log.Printf("[WS] %s\n", err)
			if !atEOF {
				return 0, nil, nil
			}
			return 0, pb, bufio.ErrFinalToken
		}
		return len(pb), pb, nil
	}
	scanner.Split(packetSplit)

	for scanner.Scan() {
		err := scanner.Err()
		if err != nil {
			if err, ok := err.(net.Error); ok && err.Timeout() {
				// socket up but silent
				log.Printf("[WS] socket up but silent %s\n", err)
			}
		}

		b := scanner.Bytes()
		p, err := packet.Start(b)
		if err != nil {
			log.Printf("[WS] Start err %s\n", err)
			return
		}
		parseErr := p.Parse()
		if parseErr != 0 {
			log.Printf("[WS] parseErr %d\n", parseErr)
		}
	}
}

func readInto(r *http.Request, c *ws.Conn, buf *bytes.Buffer) {
	for {
		err := fillBuff(r.Context(), c, buf)
		if ws.CloseStatus(err) == ws.StatusNormalClosure {
			return
		}
		if err != nil {
			log.Printf("failed to echo with %v: %v", r.RemoteAddr, err)
			return
		}
	}
}

func fillBuff(ctx context.Context, c *ws.Conn, buf *bytes.Buffer) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*100)
	defer cancel()

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}
	log.Println("[WS] message type", typ)

	_, err = io.Copy(buf, r)
	if err != nil {
		return fmt.Errorf("[WS] failed to io.Copy: %w", err)
	}
	log.Println("[WS] copied")

	return err
}
