package model

const (
	_ int = iota
	MatchOptRet_Success			// 成功
	MatchOptRet_NoMatchData		// 没有该比赛信息
	MatchOptRet_GoldNotEnough	// 金币不足
	MatchOptRet_NoPlayer		// 不存在该用户
	MatchOptRet_DelError		// 扣钱失败
)

const (
	_ int = iota
	ChangeCoinType_RedPacket
	ChangeCoinType_OneHundred
	ChangeCoinType_RandMatch
	ChangeCoinType_FriendMatch
	ChangeCoinType_Advent
	ChangeCoinType_MatchTicket
	ChangeCoinType_InviteFriend
)

type RpcCommonEmptyRequest struct {
}

type RpcCommonEmptyResponse struct{}

type RpcCommonResultResponse struct {
	Result int
}

type RpcLoadAccountInfoRequest struct {
	OpenId string
	GameId string
	SessionKey string
	NeedWxData string
	InviteOpenid string
}

type RpcLoadAccountInfoResponse struct {
	AccountInfoData AccountInfo
	UserData WxUserInfo
}

type RpcUpdateWxInfoRequest struct {
	GameId string
	WxData WxUserInfo
}

type RpcGetOneMatchInfoRequest struct {
	Gameid string
	Openid string
}

type RpcGetOneMatchInfoResponse struct {
	Result int
	CurGold uint64
	Data string
	AccData AccountInfo
	WxData WxUserInfo
}

type RpcDelMatchTicketRequest struct{
	Gameid string
	Openid string
	MatchType int
}

type RpcDelMatchTicketResponse struct {
	Result int
	CurNum uint64
}

type RpcEnrollMatchResponse struct {
	RoomId string
}

type RpcCalMatchRewardRequest struct {
	MatchType int
	Rank int
	Gameid string
	Openid string
}

type RpcCalMatchRewardResponse struct {
	Result int
	CurGold uint64
	CurMoney uint64
}

type RpcAddCoinsRequest struct{
	Gameid string
	Openid string
	AddType int
	Data []OneReward
}

type RpcAddCoinsResponse struct {
	Result int
	AccData AccountInfo
}

type RpcUpdateQDInfoRequest struct {
	Gameid string
	ChanId string
	Num string
}

type RpcGetSignInfoRequest struct {
	SignType int
	GameId string
	OpenId string
}

type RpcGetSignInfoResponse struct {
	SignSumDay int					// 累计签到天数
	IsTodaySign bool				// 今天是否签到
	RewardDatas [][]OneReward		// 签到奖励信息
}