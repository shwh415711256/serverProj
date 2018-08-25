package internal

import (
	"reflect"
	"server/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"GameServer/module/dbclient"
	"golang.org/x/net/context"
	"Common/model"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init(){
	handler(&msg.LoginReq{}, handlerLogin)
}

func handlerLogin(args []interface{}) {
	AgentMapLock.Lock()
	defer AgentMapLock.Unlock()
	m := args[0].(*msg.LoginReq)
	agent := args[1].(gate.Agent)

	log.Debug("login gameid:%s, openid:%s", m.GameId, m.OpenId)

	AgentMap[m.GameId + m.OpenId] = agent

	data := &model.AgentUserData{
		GameId:m.GameId,
		UserData:model.WxUserInfo{
			OpenId:m.OpenId,
		},
	}

	req := &model.RpcLoadAccountInfoRequest{
		OpenId: m.OpenId,
		GameId: m.GameId,
		NeedWxData: "1",
	}
	reply := &model.RpcLoadAccountInfoResponse{}
	err := dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcLoadAccountInfo", req, reply)
	if err == nil {
		data.UserData = reply.UserData
	}
	agent.SetUserData(data)
	agent.WriteMsg(&msg.LoginResp{
		Result: "ok",
	})
}