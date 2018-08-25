package internal

import (
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"Common/model"
	"sync"
)

var (
	AgentMapLock sync.RWMutex
	AgentMap map[string]gate.Agent
)

func init() {
	skeleton.RegisterChanRPC("NewAgent", rpcNewAgent)
	skeleton.RegisterChanRPC("CloseAgent", rpcCloseAgent)
	AgentMap = make(map[string]gate.Agent)
}

func rpcNewAgent(args []interface{}) {
	a := args[0].(gate.Agent)
	log.Debug("connet success, %v", a)
}

func rpcCloseAgent(args []interface{}) {
	AgentMapLock.Lock()
	defer AgentMapLock.Unlock()
	a := args[0].(gate.Agent)
	userData := a.UserData().(*model.AgentUserData)
	key := userData.GameId + userData.UserData.OpenId
	delete(AgentMap, key)
}
