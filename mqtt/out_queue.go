package mqtt

import "log"

func rangeOutQueue(connections Connections, retries Retries, outQueue <-chan OutData) {
	for o := range outQueue {
		if c, ok := connections.findConn(o.clientId); ok {
			n, err := c.publish(o.packet.toByteSlice())
			if err != nil {
				log.Println("cannot write to", o.clientId, ":", err)
			}
			// if qos > 0 save for retry
			if o.packet.QoS() > 0 {
				var r Retry
				r.clientId = o.clientId
				r.packetIdentifier = o.packet.packetIdentifier
				r.qos = o.packet.QoS()
				err = retries.addRetry(r)
				if err != nil {
					log.Println("cannot save retry of", r.clientId, r.packetIdentifier)
				}
			}
			log.Println("published", n, "bytes to", o.clientId)
		} else {
			log.Println(o.clientId, "is not connected")
		}
	}
}
