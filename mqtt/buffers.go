package mqtt

import (
	"log"
)

func rangeBuffers(buffers <-chan []byte, packets chan<- Packet) {
	var p Packet
	closed := true
	for b := range buffers {
		for len(b) > 0 {
			// log.Printf("buffers: %d bytes remaining in buffer\n", len(b))
			// log.Println(b)
			if closed {
				p, err := Start(b)
				if err != nil {
					log.Println("buffers: error starting packet", err)
					p.err = err
				}
				if p.PacketComplete() {
					packets <- p
					closed = true
				} else {
					closed = false
				}
				b = b[p.PacketLength():]
			} else {
				n := p.CompletePacket(b)
				if p.PacketComplete() {
					packets <- p
					closed = true
				} else {
					closed = false
				}
				b = b[n:]
			}
		}

	}
}
