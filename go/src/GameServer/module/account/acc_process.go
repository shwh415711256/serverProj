package account

import (
	"Common/model"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/cihub/seelog"
	"io/ioutil"
	"GameServer/module/dbclient"
	"golang.org/x/net/context"
)

// 微信登陆
func ProcessWxLogin(data *model.AccountLoginRequest) (*model.AccountLoginResponse, error){
	url := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=%s",
		data.AppId, data.Secrect, data.Code, data.GrantType)
	ret, err := http.Get(url)
	defer ret.Body.Close()
	if err != nil {
		return nil, err
	}
	bodyData, err := ioutil.ReadAll(ret.Body)
	if err != nil {
		return nil, err
	}

	seelog.Debugf("wxLogin Resp:%s", string(bodyData))
	resp := &model.AccountLoginResponse{}
	err = json.Unmarshal(bodyData, resp)
	if err != nil {
		return nil, err
	}
	req := &model.RpcLoadAccountInfoRequest{
		OpenId:resp.OpenId,
		GameId:data.GameId,
		SessionKey:resp.SessionKey,
		InviteOpenid:data.InviteOpenid,
	}
	req.GameId = data.GameId
	if req.GameId != "" {
		reply := &model.RpcLoadAccountInfoResponse{}
		err = dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcLoadAccountInfo", req, reply)
		if err != nil {
			seelog.Errorf("RpcLoadAccountInfo error[%v]", err)
		} else {
			resp.AccountInfo = reply.AccountInfoData
		}
	}
	return resp, nil
}