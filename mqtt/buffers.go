package mqtt

import (
	"log"
)

func rangeBuffers(buffers <-chan []byte, packets chan<- Packet) {
	var p Packet
	for b := range buffers {
		log.Printf("read buffer of %d bytes\n", len(b))
		log.Println(b)
		for len(b) > 0 {
			if p.PacketComplete() {
				p, err := Start(b)
				if err != nil {
					p.err = err
					packets <- p
					break
				}
				b = b[p.PacketLength():]
			} else {
				n := p.CompletePacket(b)
				b = b[n:]
			}
			if p.PacketComplete() {
				packets <- p
			}
		}

	}
}
