package mqtt

import (
	"bufio"
	"errors"
	"net"

	"github.com/ilgianlu/tagyou/packet"
	"github.com/rs/zerolog/log"
)

var ErrConnectionDown = errors.New("connection down")

func NewScanner(conn net.Conn) *bufio.Scanner {
	scanner := bufio.NewScanner(conn)
	packetSplit := func(b []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(b) == 0 && atEOF {
			// socket down - closed
			return 0, b, ErrConnectionDown
		}
		pb, err := packet.ReadFromByteSlice(b)
		if err != nil {
			log.Error().Err(err).Msg("[MQTT] error reading bytes")
			if !atEOF {
				return 0, nil, nil
			}
			return 0, pb, bufio.ErrFinalToken
		}
		return len(pb), pb, nil
	}
	scanner.Split(packetSplit)
	return scanner
}
