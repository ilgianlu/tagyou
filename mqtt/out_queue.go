package mqtt

import (
	"log"
	"time"
)

func rangeOutQueue(connections Connections, retries Retries, outQueue <-chan OutData) {
	for o := range outQueue {
		if c, ok := connections.findConn(o.clientId); ok {
			n, err := c.publish(o.packet.toByteSlice())
			if err != nil {
				log.Println("cannot write to", o.clientId, ":", err)
			} else {
				log.Println("published", n, "bytes to", o.clientId)
			}
			// if qos > 0 save for retry
			if o.packet.QoS() > 0 {
				waitForAck(o.clientId, o.packet, retries)
			}
		} else {
			log.Println(o.clientId, "is not connected")
		}
	}
}

func waitForAck(clientId string, packet Packet, retries Retries) {
	var r Retry
	r.clientId = clientId
	r.packetIdentifier = packet.packetIdentifier
	r.qos = packet.QoS()
	if r.qos == 1 {
		r.ackStatus = WAIT_FOR_PUB_ACK
	} else {
		r.ackStatus = WAIT_FOR_PUB_REC
	}
	r.applicationMessage = packet.ApplicationMessage()
	r.createdAt = time.Now()
	err := retries.addRetry(r)
	if err != nil {
		log.Println("cannot save retry of", r.clientId, r.packetIdentifier)
	}
}
