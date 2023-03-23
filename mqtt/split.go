package mqtt

import (
	"bufio"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/rs/zerolog/log"
)

func packetSplit(session *model.RunningSession, events chan<- *packet.Packet) func(b []byte, atEOF bool) (int, []byte, error) {
	return func(b []byte, atEOF bool) (advance int, token []byte, err error) {
		if len(b) == 0 && atEOF {
			// socket down - closed
			if session.GetClientId() != "" {
				log.Debug().Msgf("[MQTT] (%s:%d) will due to socket down!", session.GetClientId(), session.LastConnect)
				willEvent(session, events)
				disconnectClient(session, events)
				return 0, nil, nil
			}
			return 0, b, bufio.ErrFinalToken
		}
		pb, err := packet.ReadFromByteSlice(b)
		if err != nil {
			if !atEOF {
				return 0, nil, nil
			}
			log.Debug().Msgf("[MQTT] error reading bytes - session: %s : %s", session.GetClientId(), err.Error())
			return 0, pb, bufio.ErrFinalToken
		}
		return len(pb), pb, nil
	}
}
