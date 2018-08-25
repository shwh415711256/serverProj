package model

const (
	_ int = iota
	MatchType_RedPacket			// 红包赛
	MatchType_OneHundred			// 百元赛
	MatchType_Random				// 随机匹配
	MatchType_Friend				// 好友匹配
)

const (
	_ int = iota
	CoinType_Gold					// 金币
	CoinType_Money					// 现金
)

const (
	_ int = iota
	ConsumeCoinType_RandMatchTicket			// 随机匹配门票
	ConsumeCoinType_RandMatchLose				// 随机匹配失败
)

const (
	_ int = iota
	GetCoinType_RandMatchWin					// 随机匹配获胜
	InviteFriend_Match							// 好友对战成功邀请好友
)

type BasketBallInfo struct{
	Time int64			// 时间 单位:秒
	Score int			// 分数
}

type GameHisInfo struct{
	Gameid string
	Openid string
	Data string					// json字串
}

type BasketBallHisData struct {
	Data []BasketBallInfo
}

type BasketBallOneResult struct{
	UserData WxUserInfo
	Score int
}

type BasketBallResult struct {
	Result []BasketBallOneResult
}

type EnrollMatchRequest struct {
	Gameid string
	Openid string
	MatchType int
}