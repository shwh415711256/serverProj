package gate

import (
	"server/msg"
	"server/game"
)

func init() {
	msg.Processor.SetRouter(&msg.LoginReq{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.EnterRoomRequest{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.StartGameRequest{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.UpdateGameInfo{}, game.ChanRPC)
	msg.Processor.SetRouter(&msg.HeartBeatMsg{}, game.ChanRPC)
}
