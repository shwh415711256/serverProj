package room

import (
	"github.com/gin-gonic/gin"
	"Common/model"
	"net/http"
	"encoding/json"
	"golang.org/x/net/context"
	"github.com/cihub/seelog"
	"GameServer/module/dbclient"
)

func CreateRoom(c *gin.Context){
	req := &model.CreateRoomRequest{}
	err := c.BindJSON(req)
	if err != nil {
		merr := &model.CommonFormatError{
			Type:      model.RequestParamError,
			Status:    model.CreateRoomParamErrorStatus,
			ErrorDesc: err.Error(),
		}
		bts, err := json.Marshal(merr)
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
	}
	if req.OpenId == "" || req.GameId == "" {
		if err != nil {
			if err != nil {
				merr := &model.CommonFormatError{
					Type:      model.RequestParamError,
					Status:    model.CreateRoomParamErrorStatus,
					ErrorDesc: err.Error(),
				}
				bts, err := json.Marshal(merr)
				if err == nil {
					c.Data(http.StatusBadRequest, "application/json", bts)
				}
			}
		}
	}
	resp := &model.CreateRoomResponse{}
	err = roomManager.roomXclient.Call(context.Background(), "RpcCreateRoom", req, resp)
	if err != nil {
		seelog.Errorf("RpcCreateRoom error[%v], gameid:%s, openid:%s", err, req.GameId, req.OpenId)
	}
	seelog.Debugf("RpcCreateRoom resp:%v", resp)
	bts, err := json.Marshal(resp)
	if err == nil {
		c.Data(http.StatusOK, "application/json", bts)
	}
}

func RandomMatchOne(c *gin.Context) {
	req := &model.RpcGetOneMatchInfoRequest{}
	err := c.BindJSON(req)
	if err != nil {
		seelog.Errorf("RandomMatchOne error[%v]", err)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.RequestParamError,
			Status: model.GetOneMatchInfoRpcErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}
	if req.Gameid == "" || req.Gameid == "0" || req.Openid == "" {
		seelog.Errorf("RandomMatchOne param error,param:%v", req)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.RequestParamError,
			Status: model.GetOneMatchInfoRpcErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}

	reply := &model.RpcGetOneMatchInfoResponse{}
	err = dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcGetOneMatchInfo", req, reply)
	if err != nil {
		seelog.Errorf("RpcGetOneMatchInfo error[%v]", err)
		merr := &model.CommonFormatError{
			Type:      model.InnerError,
			Status:    model.GetOneMatchInfoRpcErrorStatus,
			ErrorDesc: "RpcGetOneMatchInfo error",
		}
		bts, err := json.Marshal(merr)
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}

	enemyResult := &model.BasketBallResult{}
	err = json.Unmarshal([]byte(reply.Data), enemyResult)
	if err != nil {
		seelog.Errorf("RpcGetOneMatchInfo json unmarshal error[%v]", err)
		merr := &model.CommonFormatError{
			Type:      model.InnerError,
			Status:    model.GetOneMatchInfoRpcErrorStatus,
			ErrorDesc: "RpcGetOneMatchInfo error",
		}
		bts, err := json.Marshal(merr)
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}

	// 扣除门票
	ticketReq := &model.RpcDelMatchTicketRequest{
		Gameid: req.Gameid,
		Openid: req.Openid,
		MatchType: model.MatchType_Random,
	}
	ticketReply := &model.RpcDelMatchTicketResponse{}
	err = dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcDelMatchTicket", ticketReq, ticketReply)
	if err != nil {
		seelog.Errorf("RpcGetMatchTicket error[%v]", err)
	}
	reply.Result = ticketReply.Result
	reply.CurGold = ticketReply.CurNum
	score := 0
	l := len(enemyResult.Result)
	for i := 0; i < l; i++ {
		score += enemyResult.Result[i].Score
	}
	roomManager.AddRandMatchEnemy(req.Gameid, req.Openid, &model.RandMatchEnemy{
		Score:score,
	})
	bts, err := json.Marshal(reply)
	if err == nil {
		c.Data(http.StatusOK, "application/json", bts)
	}
}

func CommitGameHis(c *gin.Context) {
	data := &model.GameHisInfo{}
	err := c.BindJSON(data)
	if err != nil {
		seelog.Errorf("CommitGameHis error[%v]", err)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.RequestParamError,
			Status: model.CommitGameInfoErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}
	if data.Gameid == "" || data.Gameid == "0" || data.Openid == "" {
		seelog.Errorf("CommitGameHis param error,param:%v", data)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.RequestParamError,
			Status: model.CommitGameInfoErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}
	enemyData, find := roomManager.GetRandMatchEnemy(data.Gameid, data.Openid)
	if !find {
		seelog.Errorf("CommitGameHis error, no enemy data")
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.RequestParamError,
			Status: model.CommitGameInfoErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}
	score, match := CheckGameHisStruct(data.Gameid, data.Data)
	if !match {
		seelog.Errorf("CommitGameHis error, data struct not match, data:%v", data)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.RequestParamError,
			Status: model.CommitGameInfoErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}
	reply := &model.CommonEmptyResponse{}
	err = dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcInsertGameHis", data, reply)
	if err != nil {
		seelog.Errorf("RpcCommitGameHis error[%v]", err)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.InnerError,
			Status: model.CommitGameInfoRpcErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
	}
	rank := 1
	if enemyData.Score > score {
		rank = 2
	}
	calMatchReq := &model.RpcCalMatchRewardRequest{
		Gameid:data.Gameid,
		Openid:data.Openid,
		MatchType:model.MatchType_Random,
		Rank:rank,
	}
	calMatchReply := &model.RpcCalMatchRewardResponse{}
	err = dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcCalMatchReward", calMatchReq, calMatchReply)
	if err != nil{
		seelog.Errorf("RpcCalMatchReward error[%v]", err)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.InnerError,
			Status: model.CommitGameInfoRpcErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
	}
	roomManager.RemoveRandMatchEnemy(data.Gameid, data.Openid)
	bts, err := json.Marshal(calMatchReply)
	if err != nil {
		c.Data(http.StatusOK, "application/json", bts)
	}
}

func CheckGameHisStruct(gameid string, data string) (int, bool) {
	score := 0
	switch gameid{
	case model.GameId_baketBall:
		{
			datas := &model.BasketBallHisData{}
			err := json.Unmarshal([]byte(data), datas)
			if err != nil {
				return score,false
			}
			l := len(datas.Data)
			if l == 0 {
				return score,false
			}
			last := datas.Data[0]
			if last.Score < 0 || last.Score > 3 {
				return score,false
			}
			for i := 1; i < l; i ++{
				score += datas.Data[i].Score
				if datas.Data[i].Score < 0 || datas.Data[i].Score > 3 {
					return score, false
				}
			}
		}
	default:
		{
			return score, false
		}
	}
	return score, true
}

func EnrollMatch(c *gin.Context){
	data := &model.EnrollMatchRequest{}
	err := c.BindJSON(data)
	if err != nil {
		seelog.Errorf("EnrollMatch error[%v]", err)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.RequestParamError,
			Status: model.EnrollMatchErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}
	if data.Gameid == "" || data.Gameid == "0" || data.Openid == "" {
		seelog.Errorf("EnrollMatch param error,param:%v", data)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.RequestParamError,
			Status: model.EnrollMatchErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
		return
	}
	resp := &model.RpcCommonResultResponse{}
	err = GetRoomManager().roomXclient.Call(context.Background(), "RpcEnrollMatch", data, resp)
	if err != nil {
		seelog.Errorf("RpcEnrollMatch error, gameid:%s, openid:%s, matchType:%d", data.Gameid, data.Openid, data.MatchType)
		bts, err := json.Marshal(&model.CommonFormatError{
			Type:   model.InnerError,
			Status: model.EnrollMatchErrorStatus,
		})
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
	}
}