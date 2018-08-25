package account

import "Common/model"

type AccountProvider interface{
	LoadAccountInfo(gameid string, openid string)  (*[]map[string][]byte, error)
	InsertAccountInfo(gameid string, openid string) error
	LoadWxInfo(gameid string, openid string) (*[]map[string][]byte, error)
	UpdateWxInfo(gameid string, data *model.WxUserInfo) error
	UpdateAccountValue(gameid string, openid string, args ...interface{}) error
	UpdateQDInfo(gameid string, chanid string, num string) error
}