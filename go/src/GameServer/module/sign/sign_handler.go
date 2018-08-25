package sign

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"fmt"
	"Common/model"
	"GameServer/module/dbclient"
	"golang.org/x/net/context"
	"github.com/cihub/seelog"
)

func GetSignInfo(c *gin.Context){
	gameid := c.DefaultQuery("gameid", "")
	openid := c.DefaultQuery("openid", "")
	signtypeStr := c.DefaultQuery("signtype", "0")
	if gameid == "" || openid == "" {
		seelog.Errorf("GetSignInfo error[gameid or openid is empty]")
		c.Data(http.StatusBadRequest, "text/plain", []byte("GetSignInfo error[gameid or openid is empty]"))
		return
	}
	signType, err := strconv.ParseInt(signtypeStr, 10, 32)
	if err != nil {
		seelog.Errorf("GetSignInfo parse signtype error,signtype:%s", signtypeStr)
		c.Data(http.StatusBadRequest, "text/plain", []byte(fmt.Sprintf("GetSignInfo parse signtype error[%v]", err)))
		return
	}
	req := &model.RpcGetSignInfoRequest{
		SignType:int(signType),
		GameId:gameid,
		OpenId:openid,
	}
	reply := &model.RpcGetSignInfoResponse{}
	err = dbclient.GetDBManager().GetDbxclient().Call(context.Background(), "RpcGetSignInfo", req, reply)
	if err != nil {
		seelog.Errorf("RpcGetSignInfo error[%v]", err)
		c.Data(http.StatusBadRequest, "text/plain", []byte("RpcGetSignInfo error"))
		return
	}
	c.Data(http.StatusOK, "text/plain", []byte("ok"))
}