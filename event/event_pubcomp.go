package event

import (
	"github.com/rs/zerolog/log"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
)

func clientPubcomp(clientId string, packetIdentifier int, reasonCode uint8) {
	onRetryFound := func(retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_COMP {
			persistence.RetryRepository.Delete(retry)
		} else {
			log.Info().Msgf("pubcomp for invalid retry status %s %d %d", retry.ClientId, retry.PacketIdentifier, retry.AckStatus)
		}
	}

	retry, err := persistence.RetryRepository.FirstByClientIdPacketIdentifierReasonCode(clientId, packetIdentifier, reasonCode)
	if err != nil {
		log.Info().Msgf("pubcomp for invalid retry %s %d", retry.ClientId, retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}

}
