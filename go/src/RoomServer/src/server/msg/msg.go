package msg

import (
	"github.com/name5566/leaf/network/json"
	"Common/model"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&LoginReq{})
	Processor.Register(&LoginResp{})
	Processor.Register(&EnterRoomRequest{})
	Processor.Register(&EnterRoomResponse{})
	Processor.Register(&UpdateRoomUserData{})
	Processor.Register(&StartGameRequest{})
	Processor.Register(&StartGameResponse{})
	Processor.Register(&GameResultMsg{})
	Processor.Register(&UpdateGameInfo{})
	Processor.Register(&NotifyUserLeaveRoom{})
	Processor.Register(&HeartBeatMsg{})
}

type LoginReq struct {
	GameId string
	OpenId string
}

type LoginResp struct {
	Result string
	IsReConnect bool
	GameLeftTime int64			// 游戏剩余时间 单位:微秒
	ReConnectData string
}

type EnterRoomRequest struct {
	RoomId uint64
	OpenId string
	GameId string
}

type EnterRoomResponse struct {
	Result    int
	RoomState int
	CreaterOpenid  string
	Others    []model.RoomUserData
}

type UpdateRoomUserData struct {
	RoomState int
	UserData model.RoomUserData
}

type StartGameRequest struct {
	OpenId string
	RoomId uint64
}

type StartGameResponse struct{
	Result int				// 1 成功 2 不是房间创建者 3 房间状态不对 4 房間不存在
	GamePlayTime int64		// 游戏持续时间 单位:秒
}

type CommonMsgRequest struct {
}

type HeartBeatMsg struct {
}

type GameResultMsg struct {
	Gameid string
	Data string				// json字串
}

type UpdateGameInfo struct {
	Gameid string
	Openid string
	Data string				// json字串
}

type NotifyUserLeaveRoom struct {
	CreaterOpenid string
	Openids []string
}