package mqtt

import "log"

func rangeOutQueue(connections Connections, outQueue <-chan OutData) {
	for o := range outQueue {
		if c, ok := connections.findConn(o.clientId); ok {
			n, err := c.publish(o.packet.toByteSlice())
			if err != nil {
				log.Println("cannot write to", o.clientId, ":", err)
			}
			log.Println("published", n, "bytes to", o.clientId)
		} else {
			log.Println(o.clientId, "is not connected")
		}
	}
}
