package account

import (
	"github.com/gin-gonic/gin"
	"Common/model"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin/json"
	"net/http"
	"strings"
	"net/url"
	"GameServer/module/dbclient"
	"golang.org/x/net/context"
	"strconv"
)

func Login(c *gin.Context) {
	strs := strings.Split(c.Request.RequestURI, "?")
	murl, err := url.ParseQuery(strs[2])
	if err != nil {
		seelog.Errorf("%s", err)
		return
	}
	loginRequst := &model.AccountLoginRequest{
		AppId:murl.Get("appid"),
		Secrect:murl.Get("secret"),
		Code:murl.Get("js_code"),
		GrantType:murl.Get("grant_type"),
		GameId:murl.Get("game_id"),
		InviteOpenid: murl.Get("invite_openid"),
	}
	resp, err := ProcessWxLogin(loginRequst)
	if err != nil {
		merr := &model.CommonFormatError{
			Type: model.InnerError,
			Status: model.WXLoginErrorStatus,
			ErrorDesc: err.Error(),
		}
		bts, err := json.Marshal(merr)
		if err == nil {
			c.Data(http.StatusBadRequest, "application/json", bts)
		}
	}
	bts, err := json.Marshal(resp)
	if err == nil {
		c.Data(http.StatusOK, "application/json", bts)
	}
}

func UpdateWxUserInfo(c *gin.Context) {
	wxData := &model.WxUserInfo{}
	gameId := c.DefaultQuery("gameid", "0")
	if gameId != "0" {
		err := c.BindJSON(wxData)
		if err != nil {
			merr := &model.CommonFormatError{
				Type:      model.RequestParamError,
				Status:    model.WXUpdateErrorStatus,
				ErrorDesc: err.Error(),
			}
			bts, err := json.Marshal(merr)
			if err == nil {
				c.Data(http.StatusBadRequest, "application/json", bts)
			}
		}
		if wxData.OpenId == "" {
			c.Data(http.StatusOK, "text/plain", []byte("ok"))
			return
		}
		req := &model.RpcUpdateWxInfoRequest{
			GameId:gameId,
			WxData:*wxData,
		}
		reply := &model.RpcCommonEmptyResponse{}
		err = dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcUpdateWxUserInfo", req, reply)
		if err != nil {
			seelog.Errorf("RpcUpdateWxUserInfo error[%v], gameid:%s, openid:%s", req.GameId, wxData.OpenId)
		}
	}
	seelog.Debugf("updateWxUserInfo req[%v]", wxData)
	c.Data(http.StatusOK, "text/plain", []byte("ok"))
}

func GetUserAccInfo(c *gin.Context){
	gameId := c.DefaultQuery("gameid", "")
	openId := c.DefaultQuery("openid", "")
	if gameId == "" || openId == "" {
		seelog.Errorf("GetUserAccInfo param error,gameid:%s,openid:%s", gameId, openId)
		c.Data(http.StatusBadRequest, "text/plain", []byte("GetUserAccInfo param error, gameid or openid is empty"))
		return
	}
	req := &model.RpcLoadAccountInfoRequest{
		OpenId:openId,
		GameId:gameId,
	}
	reply := &model.RpcLoadAccountInfoResponse{}
	err := dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcLoadAccountInfo", req, reply)
	if err != nil {
		seelog.Errorf("RpcLoadAccountInfo error[%v]", err)
		c.Data(http.StatusBadRequest, "text/plain", []byte("GetUserAccInfo error"))
		return
	}
	bts, err := json.Marshal(reply.AccountInfoData)
	if err == nil {
		c.Data(http.StatusOK, "application/json", bts)
	}
}

func UpdateQDLoginNum(c *gin.Context){
	chanId := c.DefaultQuery("chanId", "")
	gameId := c.DefaultQuery("gameid", "")
	num := c.DefaultQuery("num", "")
	if chanId == "" || gameId == "" || num == ""{
		seelog.Errorf("UpdateQDLoginNum error, param error,chanid:%s,gameid:%s,num:%s", chanId, gameId, num)
		c.Data(http.StatusBadRequest, "text/plain", []byte("chanId or gameId or num is empty"))
		return
	}

	_, err := strconv.ParseInt(chanId, 10, 32)
	if err != nil {
		seelog.Errorf("UpdateQDLoginNum error, param error,chanid:%s,gameid:%s,num:%s", chanId, gameId, num)
		c.Data(http.StatusBadRequest, "text/plain", []byte("chanId error"))
		return
	}
	_, err = strconv.ParseInt(gameId, 10, 32)
	if err != nil {
		seelog.Errorf("UpdateQDLoginNum error, param error,chanid:%s,gameid:%s,num:%s", chanId, gameId, num)
		c.Data(http.StatusBadRequest, "text/plain", []byte("gameId error"))
		return
	}
	_, err = strconv.ParseInt(num, 10, 32)
	if err != nil {
		seelog.Errorf("UpdateQDLoginNum error, param error,chanid:%s,gameid:%s,num:%s", chanId, gameId, num)
		c.Data(http.StatusBadRequest, "text/plain", []byte("num error"))
		return
	}

	req := &model.RpcUpdateQDInfoRequest{
		Gameid:gameId,
		ChanId:chanId,
		Num:num,
	}
	reply := &model.RpcCommonEmptyResponse{}
	err = dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcUpdateQDInfo",req, reply)
	if err != nil {
		seelog.Errorf("RpcUpdateQDInfo error[%v]", err)
		c.Data(http.StatusBadRequest, "text/plain", []byte("RpcUpdateQDInfo error"))
		return
	}
	c.Data(http.StatusOK, "text/plain", []byte("ok"))
	seelog.Debugf("UpdateQDLoginNum success,chanid:%s,gameid:%s,num:%s", chanId, gameId, num)
}