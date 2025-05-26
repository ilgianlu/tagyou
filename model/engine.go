package model

type Engine interface {
	OnConnect(session *RunningSession)
	OnPing(session *RunningSession)
	OnSubscribe(session *RunningSession, p Packet)
	OnUnsubscribe(session *RunningSession, p Packet)
	OnPublish(session *RunningSession, p Packet)
	OnClientPuback(session *RunningSession, p Packet)
	OnClientPubrec(session *RunningSession, p Packet)
	OnClientPubrel(session *RunningSession, p Packet)
	OnClientPubcomp(clientId string, packetIdentifier int, reasonCode uint8)
	OnClientDisconnect(session *RunningSession, clientId string)
	// connection status
	OnSocketUpButSilent(session *RunningSession) bool
	OnSocketDownClosed(session *RunningSession) bool
}
