package engine

import (
	"github.com/ilgianlu/tagyou/model"
	"github.com/ilgianlu/tagyou/packet"
)

func (s StandardEngine) OnPing(session *model.RunningSession) {
	toSend := packet.PingResp()
	bs, err := toSend.ToByteSlice()
	if err != nil {
		return
	}
	session.Router.Send(session.GetClientId(), bs)
}
