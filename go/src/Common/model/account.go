package model

type AccountInfo struct {
	// openid
	WxOpenId string			`json:"openid"`
	// 金币
	Gold uint64				`json:"gold"`
	// 现金 分为单位
	Money uint64			`json:"money"`
	// 分数
	Score uint64			`json:"score"`
}

type WxAccessTokenInfo struct {
	AccessToken string 		`json:"access_token"`
	ExpiresIn string		`json:"expires_in"`
}

type WxUserInfo struct {
	OpenId string			`json:"openId"`
	NickName string			`json:"nickName"`
	AvatarUrl string		`json:"avatarUrl"`
	Gendar int				`json:"gender"`
	City string				`json:"city"`
	Province string			`json:"province"`
	Country	string			`json:"country"`
	Lang string				`json:"language"`
}

type AccountLoginRequest struct {
	LoginType int 			`json:"login_type"`
	AppId string			`json:"appid"`
	Secrect string			`json:"secret"`
	Code string				`json:"code"`
	GrantType string		`json:"grant_type"`
	GameId string			`json:"game_id"`
	InviteOpenid string		`json:"invite_openid"`
}

type AccountLoginResponse struct {
	OpenId string			`json:"openid"`
	SessionKey string		`json:"session_key"`
	UnionId string			`json:"unionid"`
	AccountInfo AccountInfo `json:"account_info"`
}

type AccountTestRequest struct {
	Name string
}

type AccountTestResponse struct {
	Name string
}

type AgentUserData struct {
	GameId string
	AccInfo AccountInfo
	UserData WxUserInfo
}

type RoomUserData struct{
	State int
	AccInfo AccountInfo
	WxInfo WxUserInfo
}

