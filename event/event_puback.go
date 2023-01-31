package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func clientPuback(p *packet.Packet) {
	onRetryFound := func(retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_ACK {
			persistence.RetryRepository.Delete(retry)
		} else {
			log.Info().Msgf("puback for invalid retry status %s %d %d", retry.ClientId, retry.PacketIdentifier, retry.AckStatus)
		}
	}

	retry, err := persistence.RetryRepository.FirstByClientIdPacketIdentifier(p.Session.GetClientId(), p.PacketIdentifier())
	if err != nil {
		log.Info().Msgf("puback for invalid retry %s %d", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
