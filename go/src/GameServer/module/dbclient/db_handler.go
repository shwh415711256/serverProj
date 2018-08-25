package dbclient

import (
	"github.com/gin-gonic/gin"
	"Common/model"
	"context"
	"github.com/cihub/seelog"
	"net/http"
	"strconv"
)

func GetVersionIgnoreInfo(c *gin.Context){
	key := c.DefaultQuery("game_id", "1")
	req := &model.RpcVersionInfoRequest{
		key,
	}
	reply := &model.RpcVersionInfoResponse{}
	err := GetDBManager().dbxclient.Call(context.Background(), "GetVersionInfo", req, reply)
	if err != nil {
		seelog.Errorf("GetVersionIgnoreInfo error[%v]", err)
		c.Data(400, "text/plain", []byte("failed"))
		return
	}
	c.Data(http.StatusOK, "text/plain", []byte(strconv.Itoa(reply.IgnoreFlag)))
}

func ReloadVersionIgnoreInfo(c *gin.Context){
	req := &model.CommonEmptyRequest{
	}
	reply := &model.CommonEmptyResponse{}
	err := GetDBManager().dbxclient.Call(context.Background(), "ReloadVersionInfo", req, reply)
	if err != nil {
		seelog.Errorf("GetVersionIgnoreInfo error[%v]", err)
		c.Data(400, "text/plain", []byte("failed"))
		return
	}
	c.Data(http.StatusOK, "text/plain", []byte("ok"))
}
