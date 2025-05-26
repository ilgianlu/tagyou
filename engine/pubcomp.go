package engine

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
)

func (s StandardEngine) OnClientPubcomp(clientId string, packetIdentifier int, reasonCode uint8) {
	onRetryFound := func(retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_COMP {
			persistence.RetryRepository.Delete(retry)
		} else {
			slog.Info("pubcomp for invalid retry status", "client-id", retry.ClientId, "packet-identifier", retry.PacketIdentifier, "ack-status", retry.AckStatus)
		}
	}

	retry, err := persistence.RetryRepository.FirstByClientIdPacketIdentifierReasonCode(clientId, packetIdentifier, reasonCode)
	if err != nil {
		slog.Info("pubcomp for invalid retry", "client-id", retry.ClientId, "packet-identifier", retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}

}
