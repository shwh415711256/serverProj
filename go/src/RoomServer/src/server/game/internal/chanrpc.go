package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"Common/model"
	"server/game/internal/room"
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Debug("connet success, %v", a)
}

func rpcCloseAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	if a.UserData() != nil {
		userData := a.UserData().(*model.AgentUserData)
		room.LeaveRoom(userData.GameId, userData.UserData.OpenId)
	}
	log.Debug("connct close, %v", a)
}
