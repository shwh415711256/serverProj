package room

import (
	"Common/model"
	"sync"
	"server/conf"
	"errors"
	"fmt"
	"time"
	"github.com/name5566/leaf/log"
	"server/msg"
	"github.com/name5566/leaf/gate"
	"encoding/json"
	"server/game/internal/rpcclient"
	"golang.org/x/net/context"
)

// 玩家游戏状态
const (
	_ int = iota
	PlayerState_Wait 	  // 等待状态
	PlayerState_Ready	 //	 准备状态
	PlayerState_Playing	 //  游戏进行状态
	PlayerState_Out		// 退出状态
)

// 房间状态
const(
	_ int = iota
	RoomState_Ready			// 准备
	RoomState_WaitStart		// 等待开始
	RoomState_Start			// 游戏进行中
	RoomState_WaitClose		// 等待关闭
)

// 房间类型
const(
	_ int = iota
	RoomType_Friend			// 好友房
)

// 进入房间结果
const (
	_ int = iota
	EnterRoomRet_Success		// 1 进入房间成功
	EnterRoomRet_NoRoom		// 2 房间不存在
	EnterRoomRet_Full			// 3 房间已满
	EnterRoomRet_InGame		// 4 房间已经开始游戏
	EnterRoomRet_HaveInRoom 		// 5 已经在房间
)

// 开始游戏结果
const (
	_ int = iota
	StartGame_Success			// 1 成功
	StartGame_NoCreater		// 2 不是创建者
	StartGame_NoWaitStart		// 3 房间状态不是等待开始
	StartGame_NoRoom			// 4 没有该房间
)

type RoomPlayerData struct {
	GameState int
	Score int
	Agent gate.Agent
	AccInfo model.AccountInfo
	WxInfo model.WxUserInfo
}

type RoomData struct {
	// 房间id
	RoomId uint64
	// 游戏id
	RoomGameId string
	// 房间状态
	RoomState int
	// 房间基本信息
	RoomData *model.OneRoomInfo
	// 创建者openid
	CreateUserId string
	// 创建时间
	CreateTime int64
	// 游戏开始时间
	GameStartTime int64
//	UserMapLock sync.RWMutex
	UserMap map[string]*RoomPlayerData
}

type RoomManagerStruct struct {
	Lock sync.RWMutex
	RoomMap map[uint64]*RoomData
	UserRoomMap map[string]uint64
	InitId uint64
	StartRoomList []RoomData
	RecycleRoomList []uint64
}

var (
	RoomManager *RoomManagerStruct
)

func Init() {
	RoomManager = &RoomManagerStruct{
		RoomMap : make(map[uint64]*RoomData),
		UserRoomMap : make(map[string]uint64),
		StartRoomList : make([]RoomData, 0, 10),
		RecycleRoomList : make([]uint64, 0, 10),
		InitId : 1,
	}
}

func UpdateOneSec(){
	RoomManager.Lock.Lock()
	defer RoomManager.Lock.Unlock()
	now := time.Now()
	l := len(RoomManager.StartRoomList)
	for i := 0; i < l; i ++ {
		room := RoomManager.StartRoomList[i]
		playTime := now.UnixNano() - room.GameStartTime
		log.Debug("playTime:%d, startTime:%d", playTime, room.GameStartTime)
		if playTime >= room.RoomData.RoomPlayGameTime * 1000000000{
			// 游戏结束，结算
			if GameOver(room.RoomId, i) {
				l --
				i --
			}
		}
	}
	rl := len(RoomManager.RecycleRoomList)
	for i := 0; i < rl; i ++ {
		delete(RoomManager.RoomMap, RoomManager.RecycleRoomList[i])
		log.Debug("Room Recycle success,roomid:%s", RoomManager.RecycleRoomList[i])
	}
	RoomManager.RecycleRoomList = []uint64{}
}

func CheckRoomCanStart(roomData *RoomData) bool{
	if roomData.RoomState != RoomState_Ready {
		return false
	}
//	roomData.UserMapLock.RLock()
//	defer roomData.UserMapLock.RUnlock()
	l := len(roomData.UserMap)
	if l < roomData.RoomData.RoomNeedPlayerNum {
		return false
	}
	for _, u := range roomData.UserMap {
		if u.GameState != PlayerState_Ready {
			log.Debug("openid:%s", u.WxInfo.OpenId)
			return false
		}
	}
	roomData.RoomState = RoomState_WaitStart
	log.Debug("room waitStart, roomid:%s", roomData.RoomId)
	return true
}

func BroadCastMeToOthersOnEnterRoom(roomData *RoomData, userData *RoomPlayerData) {
//	roomData.UserMapLock.RLock()
//	defer roomData.UserMapLock.RUnlock()
	for _, u := range roomData.UserMap {
		if u.Agent != nil && u.AccInfo.WxOpenId != userData.AccInfo.WxOpenId{
			sendMsg := &msg.UpdateRoomUserData{
				RoomState: roomData.RoomState,
				UserData: model.RoomUserData{
					State:   userData.GameState,
					AccInfo: userData.AccInfo,
					WxInfo:  userData.WxInfo,
				},
			}
			u.Agent.WriteMsg(sendMsg)
			log.Debug("%v", sendMsg)
		}
	}
}

func CreateOneRoom(data *model.CreateRoomRequest, resp *model.CreateRoomResponse) error{
	RoomManager.Lock.Lock()
	defer RoomManager.Lock.Unlock()
	key := data.GameId + data.OpenId
	if rid, find := RoomManager.UserRoomMap[key]; find {
		return errors.New(fmt.Sprintf("user have in one room,gameid:%s, openid:%s, roomid:%s", data.GameId, data.OpenId, rid))
	}
	confRoomData := conf.GetConfRoomInfo(data.GameId, data.RoomType)
	if confRoomData == nil {
		return errors.New(fmt.Sprintf("conf have no roomdata, gameid:%s, roomtype:%s", data.GameId, data.RoomType))
	}
	room := &RoomData{
		RoomId:RoomManager.InitId,
		RoomGameId:data.GameId,
		RoomState:RoomState_Ready,
		RoomData:confRoomData,
		CreateUserId:data.OpenId,
		CreateTime:time.Now().UnixNano(),
		GameStartTime:0,
		UserMap:make(map[string]*RoomPlayerData),
	}
	playerData := &RoomPlayerData{
	}
	if room.RoomData.NeedPlayerReady == 0 {
		playerData.GameState = PlayerState_Ready
	}else{
		playerData.GameState = PlayerState_Wait
	}
//	room.UserMapLock.Lock()
	room.UserMap[data.OpenId] = playerData
//	room.UserMapLock.Unlock()
	RoomManager.RoomMap[room.RoomId] = room
	RoomManager.UserRoomMap[data.GameId + data.OpenId] = room.RoomId
	RoomManager.InitId ++
	log.Debug("create room success,roomid:%s,creatorid:%s, gameid:%s", room.RoomId, room.CreateUserId, data.GameId)
	resp.Result = 1
	resp.RoomId = room.RoomId
	resp.RoomServerAddr = "wss://cc.h5youxi.fun/wss"
	return nil
}

func EnterRoom(data *msg.EnterRoomRequest, resp *msg.EnterRoomResponse, userData *model.AgentUserData, agent gate.Agent){
	RoomManager.Lock.RLock()
	defer RoomManager.Lock.RUnlock()
	room, find := RoomManager.RoomMap[data.RoomId]
	if !find {
		resp.Result = EnterRoomRet_NoRoom
		log.Error("room not exist,roomid:%s", data.RoomId)
		return
	}

	if room.RoomState == RoomState_WaitClose {
		resp.Result = EnterRoomRet_NoRoom
		log.Error("room not exist,roomid:%s", room.RoomId)
		return
	}else if room.RoomState == RoomState_Start {
		log.Error("room game start,roomid:%s", room.RoomId)
		resp.Result = EnterRoomRet_InGame
		return
	}
//	room.UserMapLock.RLock()
	l := len(room.UserMap)
	if l >= room.RoomData.RoomMaxNum {
		resp.Result = EnterRoomRet_Full
		return
	}
	if _, find := room.UserMap[data.OpenId]; find{
		resp.Result = EnterRoomRet_HaveInRoom
		return
	}

	// 如果进来时没人，设置房间创建者
	if l == 0 {
		room.CreateUserId = data.OpenId
	}
	resp.CreaterOpenid = room.CreateUserId
	for _, u := range room.UserMap {
		resp.Others = append(resp.Others, model.RoomUserData{
			State:u.GameState,
			AccInfo:u.AccInfo,
			WxInfo:u.WxInfo,
		})
	}
//	room.UserMapLock.RUnlock()
	playerData := &RoomPlayerData{
		Agent:agent,
		AccInfo:userData.AccInfo,
		WxInfo:userData.UserData,
	}
	if room.RoomData.NeedPlayerReady == 0 {
		playerData.GameState = PlayerState_Ready
	}else{
		playerData.GameState = PlayerState_Wait
	}
//	room.UserMapLock.Lock()
	room.UserMap[data.OpenId] = playerData
	if playerData.GameState == PlayerState_Ready {
		CheckRoomCanStart(room)
	}
//	room.UserMapLock.Unlock()
	resp.Result = EnterRoomRet_Success
	resp.RoomState = room.RoomState
	RoomManager.UserRoomMap[data.GameId + data.OpenId] = room.RoomId

	// broadcast me to others
	BroadCastMeToOthersOnEnterRoom(room, playerData)

//	room.UserMapLock.RLock()
	for _, u := range room.UserMap{
		log.Debug("房间内玩家:roomid:%d, userData:%v", room.RoomId, u)
	}
	log.Debug("enter room success,roomid:%s, gameid:%s, openid:%s", room.RoomId, data.GameId, data.OpenId)
//	room.UserMapLock.RUnlock()
}

func StartGame(req *msg.StartGameRequest, resp *msg.StartGameResponse) {
	RoomManager.Lock.Lock()
	defer RoomManager.Lock.Unlock()
	room, find := RoomManager.RoomMap[req.RoomId]
	if !find {
		resp.Result = StartGame_NoRoom
		return
	}
	if room.RoomState != RoomState_WaitStart {
		resp.Result = StartGame_NoWaitStart
		log.Debug("StartGame fail, roomstate:%d", room.RoomState)
		return
	}
	if req.OpenId != room.CreateUserId {
		resp.Result = StartGame_NoCreater
		return
	}
	room.RoomState = RoomState_Start
	room.GameStartTime = time.Now().UnixNano()
	RoomManager.StartRoomList = append(RoomManager.StartRoomList, *room)
	resp.Result = StartGame_Success
	resp.GamePlayTime = room.RoomData.RoomPlayGameTime
	for _, u := range room.UserMap {
		if u.Agent != nil{
			u.Agent.WriteMsg(resp)
		}
		u.GameState = PlayerState_Playing
		if u.AccInfo.WxOpenId == room.CreateUserId {
			// 邀请好友对战加金币
			req := &model.RpcAddCoinsRequest{
				Gameid:room.RoomGameId,
				Openid:u.AccInfo.WxOpenId,
				Data:[]model.OneReward{model.OneReward{RewardType:model.CoinType_Gold, Num:10}},
			}
			reply := &model.RpcAddCoinsResponse{}
			err := rpcclient.DbXclient.Call(context.Background(), "RpcAddCoins", req, reply)
			if err != nil || reply.Result == 0{
				log.Error("RpcAddCoins error")
			}
			u.AccInfo = reply.AccData
		}
	}
	log.Debug("StartGame success, roomstate:%d", room.RoomState)
}

func IsCreaterThenSetAgent(gameId string, openId string, agent gate.Agent) {
	RoomManager.Lock.Lock()
	defer RoomManager.Lock.Unlock()
	key := gameId + openId
	rid, find := RoomManager.UserRoomMap[key]
	if !find {
		return
	}
	r, find := RoomManager.RoomMap[rid]
	if !find {
		return
	}
	u, find := r.UserMap[openId]
	if !find {
		return
	}
	u.Agent = agent
	uData := agent.UserData().(*model.AgentUserData)
	u.AccInfo = uData.AccInfo
	u.WxInfo = uData.UserData
}

func IsUserConnect(gameId string, openId string) bool{
	RoomManager.Lock.RLock()
	defer RoomManager.Lock.RUnlock()
	key := gameId + openId
	rid, find := RoomManager.UserRoomMap[key]
	if !find {
		return false
	}
	r, find := RoomManager.RoomMap[rid]
	if !find {
		return false
	}
	u, find := r.UserMap[openId]
	if !find {
		return false
	}
	if r.CreateUserId == openId {
		return false
	}
	return u.GameState == PlayerState_Playing
}

func OnUserReconnect(gameId string, openId string, agent gate.Agent, resp *msg.LoginResp) (*model.WxUserInfo, *model.AccountInfo){
	RoomManager.Lock.Lock()
	defer RoomManager.Lock.Unlock()
	key := gameId + openId
	rid, find := RoomManager.UserRoomMap[key]
	if !find {
		return nil, nil
	}
	r, find := RoomManager.RoomMap[rid]
	if !find {
		return nil, nil
	}
	u, find := r.UserMap[openId]
	if !find {
		return nil, nil
	}
	u.Agent = agent
	resp.GameLeftTime = (time.Now().UnixNano() - r.GameStartTime) / 1000000
	// 消息推送
	switch gameId {
	case model.GameId_baketBall:
		{
			info := &model.BasketBallResult{}
			for _, u := range r.UserMap {
				info.Result = append(info.Result, model.BasketBallOneResult{
					UserData:u.WxInfo,
					Score:u.Score,
				})
			}
			bts, _ := json.Marshal(info)
			resp.ReConnectData = string(bts)
		}
	}
	return &u.WxInfo, &u.AccInfo
}

func GameOver(roomId uint64, rindex int) bool{
	room, find := RoomManager.RoomMap[roomId]
	if !find {
		log.Error("GameOverResult error, can not find room, roomid:%ul", roomId)
		return false
	}
	result := &msg.GameResultMsg{}
	data := &model.BasketBallResult{}
	for _, u := range room.UserMap{
		data.Result = append(data.Result, model.BasketBallOneResult{
			UserData:u.WxInfo,
			Score:u.Score,
		})
	}
	bts, err := json.Marshal(data)
	if err != nil {
		log.Error("GameOverResult json marshal error[%v]", err)
		return false
	}
	result.Data = string(bts)
	m := &msg.NotifyUserLeaveRoom{
		CreaterOpenid:room.CreateUserId,
	}
	for _, u := range room.UserMap{
		if u.Agent != nil {
			u.Agent.WriteMsg(result)
		}else{
			m.Openids = append(m.Openids, u.WxInfo.OpenId)
		}
	}
	l := len(m.Openids)
	for i := 0; i < l; i ++ {
		delete(room.UserMap, m.Openids[i])
		delete(RoomManager.UserRoomMap, m.Openids[i])
	}
	for _, u := range room.UserMap{
		if u.Agent != nil && l > 0 {
			u.Agent.WriteMsg(m)
		}
		if room.RoomData.NeedPlayerReady == 0 {
			u.GameState = PlayerState_Ready
		}else{
			u.GameState = PlayerState_Wait
		}
	}
	if room.RoomData.NeedPlayerReady == 0 && len(room.UserMap) == room.RoomData.RoomNeedPlayerNum {
		room.RoomState = RoomState_WaitStart
	}else{
		room.RoomState = RoomState_Ready
	}
	log.Debug("GameOver Send Result Success, roomid:%d", room.RoomId)
	l = len(RoomManager.StartRoomList)
	if l == 1 {
		RoomManager.StartRoomList = []RoomData{}
	}else if rindex == l - 1 {
		RoomManager.StartRoomList = append(RoomManager.StartRoomList[:l - 1])
	}else {
		RoomManager.StartRoomList = append(RoomManager.StartRoomList[:rindex], RoomManager.StartRoomList[rindex+1:]...)
	}
	CheckRoomRecycle(room)
	return true
}

func LeaveRoom(gameId string, openId string) {
	RoomManager.Lock.Lock()
	defer RoomManager.Lock.Unlock()
	key := gameId + openId
	rid, find := RoomManager.UserRoomMap[key]
	if !find {
		log.Error("LeaveRoom error, not in a room, gameid:%s, openid:%s", gameId, openId)
		return
	}
	room, find := RoomManager.RoomMap[rid]
	if !find {
		log.Error("LeaveRoom error, room not exist, gameid:%s, openid:%s, roomid:%d", gameId, openId, rid)
		return
	}
	m := &msg.NotifyUserLeaveRoom{}
	u,_ := room.UserMap[openId]
	if u.GameState == PlayerState_Playing {
		u.Agent = nil
		return
	}
	delete(room.UserMap, openId)
	delete(RoomManager.UserRoomMap, key)
	l := len(room.UserMap)
	if l == 0 {
		room.CreateUserId = ""
	}else {
		m.CreaterOpenid = room.CreateUserId
		m.Openids = append(m.Openids, openId)
		for _, u := range room.UserMap {
			if openId == room.CreateUserId {
				room.CreateUserId = u.AccInfo.WxOpenId
				m.CreaterOpenid = room.CreateUserId
			}
			if u.Agent != nil {
				u.Agent.WriteMsg(m)
			}
		}
	}
	ready := 0
	if room.RoomState == RoomState_WaitStart {
		for _, u := range room.UserMap{
			if u.GameState == PlayerState_Ready {
				ready ++
			}
		}
		if ready < room.RoomData.RoomNeedPlayerNum {
			room.RoomState = RoomState_Ready
		}
	}
	// 广播离开
	log.Debug("LeaveRoom success,gameid:%s, openid:%s", gameId, openId)
	CheckRoomRecycle(room)
}

func CheckRoomRecycle(room *RoomData) bool{
	if room.RoomState != RoomState_Ready {
		return false
	}
	if len(room.UserMap) == 0 {
		room.RoomState = RoomState_WaitClose
		RoomManager.RecycleRoomList = append(RoomManager.RecycleRoomList, room.RoomId)
		log.Debug("Room Enter Recycle,roomid:%s", room.RoomId)
	}
	return true
}

func DealGameInfo(userData *RoomPlayerData, gameid string, data string) bool{
	switch gameid{
	case model.GameId_baketBall:
		{
			info := &model.BasketBallInfo{}
			err := json.Unmarshal([]byte(data), info)
			if err != nil{
				log.Error("DealGameInfo error, gameid:%s, data:%s", gameid, data)
				return false
			}
			if info.Score < 0 || info.Score > 3 {
				log.Error("DealGameInfo score error, gameid:%s, data:%s, uScore:%d", gameid, data, userData.Score)
				return false
			}
			userData.Score += info.Score
			return true
		}
	default:
		return false
	}
	return false
}

func UpdateUserGameInfo(data *msg.UpdateGameInfo) {
	RoomManager.Lock.RLock()
	defer RoomManager.Lock.RUnlock()
	key := data.Gameid + data.Openid
	rid, find := RoomManager.UserRoomMap[key]
	if !find {
		log.Error("UpdateUserGameInfo error, not in a room,gameid:%s, openid:%s", data.Gameid, data.Openid)
		return
	}
	room, find := RoomManager.RoomMap[rid]
	if !find {
		log.Error("UpdateUserGameInfo error, room not exist,gameid:%s, openid:%s, roomid:%s", data.Gameid, data.Openid, rid)
		return
	}
	if room.RoomState != RoomState_Start {
		log.Error("UpdateUserGameInfo error, room is not in game, roomid:%s", room.RoomId)
		return
	}
	user, find := room.UserMap[data.Openid]
	if !find {
		log.Error("UpdateUserGameInfo error, user not exist,gameid:%s, openid:%s", data.Gameid, data.Openid)
		return
	}
	if !DealGameInfo(user, data.Gameid, data.Data) {
		log.Error("DealGameInfo error")
		return
	}
	for _, u := range room.UserMap {
		if /*u.WxInfo.OpenId != data.Openid &&*/ u.Agent != nil {
			u.Agent.WriteMsg(data)
		}
	}
}
