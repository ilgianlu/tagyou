package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
	"github.com/ilgianlu/tagyou/sender"
)

func clientPubrel(sender sender.Sender, p *packet.Packet) {
	sendPubcomp := func(retry model.Retry) {
		toSend := packet.Pubcomp(p.PacketIdentifier(), retry.ReasonCode, p.Session.ProtocolVersion)
		sender.Send(p.Session.GetClientId(), toSend.ToByteSlice())
	}

	onExpectedPubrel := func(retry model.Retry) {
		sendPubcomp(retry)
		persistence.RetryRepository.Delete(retry)
	}

	onRetryFound := func(retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_REL {
			onExpectedPubrel(retry)
		} else {
			log.Info().Msgf("pubrel for invalid retry status %s %d %d", retry.ClientId, retry.PacketIdentifier, retry.AckStatus)
		}
	}
	retry, err := persistence.RetryRepository.FirstByClientIdPacketIdentifierReasonCode(p.Session.GetClientId(), p.PacketIdentifier(), p.ReasonCode)
	if err != nil {
		log.Info().Msgf("pubrel for invalid retry %s %d", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
