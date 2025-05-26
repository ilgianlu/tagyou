package event

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/persistence"
)

func OnClientPuback(session *model.RunningSession, p model.Packet) {
	onRetryFound := func(retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_ACK {
			persistence.RetryRepository.Delete(retry)
		} else {
			slog.Info("puback for invalid retry status", "client-id", retry.ClientId, "packet-identifier", retry.PacketIdentifier, "ack-status", retry.AckStatus)
		}
	}

	retry, err := persistence.RetryRepository.FirstByClientIdPacketIdentifier(session.GetClientId(), p.PacketIdentifier())
	if err != nil {
		slog.Info("puback for invalid retry", "client-id", retry.ClientId, "packet-identifier", retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
