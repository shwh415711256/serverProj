package msg

import (
	"github.com/name5566/leaf/network/json"
)

var Processor = json.NewProcessor()

func init() {
	Processor.Register(&LoginReq{})
	Processor.Register(&LoginResp{})
}

type LoginReq struct {
	GameId string
	OpenId string
}

type LoginResp struct {
	Result string
}

type CommonMsgRequest struct {

}
