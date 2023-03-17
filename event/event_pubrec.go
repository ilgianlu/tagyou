package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/routers"
)

func clientPubrec(router routers.Router, p *packet.Packet) {
	sendPubrel := func(retry model.Retry) {
		toSend := packet.Pubrel(retry.PacketIdentifier, retry.ReasonCode, p.Session.ProtocolVersion)
		router.Send(retry.ClientId, toSend.ToByteSlice())
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
