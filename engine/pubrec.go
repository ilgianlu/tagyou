package engine

import (
	"log/slog"

	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
	"github.com/ilgianlu/tagyou/persistence"
)

func (s StandardEngine) OnClientPubrec(session *model.RunningSession, p model.Packet) {
	sendPubrel := func(retry model.Retry) {
		toSend := packet.Pubrel(retry.PacketIdentifier, retry.ReasonCode, session.ProtocolVersion)
		session.Router.Send(retry.ClientId, toSend.ToByteSlice())
	}

	onExpectedPubrec := func(retry model.Retry) {
		sendPubrel(retry)
		// change retry state to wait for pubcomp
		persistence.RetryRepository.UpdateAckStatus(retry.ID, model.WAIT_FOR_PUB_COMP)
	}

	onRetryFound := func(retry model.Retry) {
		// if retry in wait for pub rec -> send pub rel
		if retry.AckStatus == model.WAIT_FOR_PUB_REC {
			onExpectedPubrec(retry)
		} else {
			slog.Info("pubrec for invalid retry status", "client-id", retry.ClientId, "packet-identifier", retry.PacketIdentifier, "ack-status", retry.AckStatus)
		}
	}

	retry, err := persistence.RetryRepository.FirstByClientIdPacketIdentifierReasonCode(session.GetClientId(), p.PacketIdentifier(), p.GetReasonCode())
	if err != nil {
		slog.Info("pubrec for invalid retry", "client-id", retry.ClientId, "packet-identifier", retry.PacketIdentifier)
	} else {
		onRetryFound(retry)
	}
}
