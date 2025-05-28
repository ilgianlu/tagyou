package engine

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func (s StandardEngine) OnClientPubrel(session *model.RunningSession, p model.Packet) {
	sendPubcomp := func(retry model.Retry) {
		toSend := packet.Pubcomp(p.PacketIdentifier(), retry.ReasonCode, session.ProtocolVersion)
		session.Router.Send(session.GetClientId(), toSend.ToByteSlice())
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
			slog.Info("pubrel for invalid retry status", "client-id", retry.ClientId, "packet-identifier", retry.PacketIdentifier, "ack-status", retry.AckStatus)
		}
	}
	retry, err := persistence.RetryRepository.FirstByClientIdPacketIdentifierReasonCode(session.GetClientId(), p.PacketIdentifier(), p.GetReasonCode())
	if err != nil {
		slog.Info("pubrel for invalid retry", "client-id", retry.ClientId, "packet-identifier", retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
