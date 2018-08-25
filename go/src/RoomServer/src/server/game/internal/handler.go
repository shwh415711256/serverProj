package internal

import (
	"reflect"
	"server/msg"
	"github.com/name5566/leaf/gate"
	"github.com/name5566/leaf/log"
	"Common/model"
	"golang.org/x/net/context"
	"server/game/internal/room"
	"server/game/internal/rpcclient"
)

func handler(m interface{}, h interface{}) {
	skeleton.RegisterChanRPC(reflect.TypeOf(m), h)
}

func init(){
	handler(&msg.LoginReq{}, handlerLogin)
	handler(&msg.EnterRoomRequest{}, handlerEnterRoom)
	handler(&msg.StartGameRequest{}, handlerStartGame)
	handler(&msg.UpdateGameInfo{}, handlerUpdateGameInfo)
	handler(&msg.HeartBeatMsg{}, handlerHeartBeat)
}

func handlerLogin(args []interface{}) {
	m := args[0].(*msg.LoginReq)
	agent := args[1].(gate.Agent)

	log.Debug("login gameid:%s, openid:%s", m.GameId, m.OpenId)

	data := &model.AgentUserData{
		GameId:m.GameId,
		UserData:model.WxUserInfo{
			OpenId:m.OpenId,
		},
	}

	reConnect := false
	if room.IsUserConnect(m.GameId, m.OpenId) {
		reConnect = true
	}

	resp := &msg.LoginResp{
		Result: "ok",
		IsReConnect:reConnect,
	}
	if !reConnect {
		req := &model.RpcLoadAccountInfoRequest{
			OpenId:     m.OpenId,
			GameId:     m.GameId,
			NeedWxData: "1",
		}
		reply := &model.RpcLoadAccountInfoResponse{}
		err := rpcclient.DbXclient.Call(context.Background(), "RpcLoadAccountInfo", req, reply)
		if err == nil {
			data.UserData = reply.UserData
			data.AccInfo = reply.AccountInfoData
		}
		agent.SetUserData(data)

		room.IsCreaterThenSetAgent(m.GameId, m.OpenId, agent)
	}else{
		wxData, accData := room.OnUserReconnect(m.GameId, m.OpenId, agent, resp)
		data.UserData = *wxData
		data.AccInfo = *accData
		log.Debug("reconnct gameid:%s, openid:%s", m.GameId, m.OpenId)
		agent.SetUserData(data)
	}
	log.Debug("login resp:%v", resp)
	agent.WriteMsg(resp)
}

func handlerEnterRoom(args []interface{}){
	m := args[0].(*msg.EnterRoomRequest)
	a := args[1].(gate.Agent)
	userData := a.UserData().(*model.AgentUserData)
	resp := &msg.EnterRoomResponse{}
	room.EnterRoom(m, resp, userData, a)
	a.WriteMsg(resp)
}

func handlerStartGame(args []interface{}){
	m := args[0].(*msg.StartGameRequest)
	resp := &msg.StartGameResponse{}
	room.StartGame(m, resp)
}

func handlerUpdateGameInfo(args []interface{}){
	m := args[0].(*msg.UpdateGameInfo)
	room.UpdateUserGameInfo(m)
}

func handlerHeartBeat(args []interface{}){
	a := args[1].(gate.Agent)
	if a.UserData() != nil {
		userData := a.UserData().(*model.AgentUserData)
		log.Debug("heart beat gameid:%s, openid:%s", userData.GameId, userData.AccInfo.WxOpenId)
	}
}