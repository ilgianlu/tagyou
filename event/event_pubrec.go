package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/out"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func clientPubrec(p *packet.Packet, outQueue chan<- out.OutData) {
	sendPubrel := func(retry model.Retry) {
		var o out.OutData
		o.ClientId = retry.ClientId
		toSend := packet.Pubrel(retry.PacketIdentifier, retry.ReasonCode, p.Session.ProtocolVersion)
		o.Packet = toSend.ToByteSlice()
		outQueue <- o
	}

	onExpectedPubrec := func(retry model.Retry) {
		sendPubrel(retry)
		// change retry state to wait for pubcomp
		retry.AckStatus = model.WAIT_FOR_PUB_COMP
		persistence.RetryRepository.SaveOne(retry)
	}

	onRetryFound := func(retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_REC {
			onExpectedPubrec(retry)
		} else {
			log.Info().Msgf("pubrec for invalid retry status %s %d %d", retry.ClientId, retry.PacketIdentifier, retry.AckStatus)
		}
	}

	retry, err := persistence.RetryRepository.FirstByClientIdPacketIdentifierReasonCode(p.Session.GetClientId(), p.PacketIdentifier(), p.ReasonCode)
	if err != nil {
		log.Info().Msgf("pubrec for invalid retry %s %d", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
