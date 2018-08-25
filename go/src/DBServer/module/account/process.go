package account

import (
	"Common/model"
	"DBServer/store/account"
	"errors"
	"fmt"
	"strconv"
	"github.com/cihub/seelog"
)

var (
	accountProvider account.AccountProvider
)

func Init() {
	accountProvider = &account.AccountDefaultProvider{}
	accountProvider.(*account.AccountDefaultProvider).Init()
}


func LoadAccountInfoFromDB(gameId string, openid string) (*model.AccountInfo, error){
	rows, err := accountProvider.LoadAccountInfo(gameId, openid)
	if err != nil {
		return nil, err
	}
	if len(*rows) == 0 {
		return nil, errors.New(fmt.Sprintf("db have no accountdata, gameId:%s, openid:%s", gameId, openid))
	}
	r := (*rows)[0]
	info := &model.AccountInfo{
		WxOpenId: string(r["openid"]),
	}
	info.Money, err = strconv.ParseUint(string(r["money"]), 10, 64)
	if err != nil {
		return nil, errors.New("parse money error")
	}
	info.Gold, err = strconv.ParseUint(string(r["gold"]), 10, 64)
	if err != nil {
		return nil, errors.New("parse gold error")
	}
	info.Score, err = strconv.ParseUint(string(r["score"]), 10, 64)
	if err != nil {
		return nil, errors.New("parse score error")
	}
	return info, nil
}

func LoadWxUserDataFromDB(gameId string, openid string) (*model.WxUserInfo, error) {
	rows, err := accountProvider.LoadWxInfo(gameId, openid)
	if err != nil {
		return nil, err
	}
	if len(*rows) == 0 {
		return nil, errors.New(fmt.Sprintf("db have no wxdata, gameId:%s, openid:%s", gameId, openid))
	}
	r := (*rows)[0]
	gender, _ := strconv.ParseInt(string(r["gendar"]), 10, 32)
	info := &model.WxUserInfo{
		OpenId:openid,
		NickName:string(r["nick_name"]),
		AvatarUrl:string(r["avatar_url"]),
		Gendar:int(gender),
		City:string(r["city"]),
		Province:string(r["province"]),
		Country:string(r["country"]),
		Lang:string(r["language"]),
	}
	return info, nil
}

func NewAccount(gameid string, openid string) *model.AccountInfo{
	accountProvider.InsertAccountInfo(gameid, openid)
	return nil
}

func UpdateWxInfo(gameid string, data *model.WxUserInfo) error{
	return accountProvider.UpdateWxInfo(gameid, data)
}

func CalculateCoins(info *model.AccountInfo, reward *model.OneReward){
	switch reward.RewardType{
	case model.CoinType_Gold:
		{
			info.Gold += uint64(reward.Num)
		}
	case model.CoinType_Money:
		{
			info.Money += uint64(reward.Num)
		}
	}
}

func UpdateCoins(gameid string, openid string, datas []model.OneReward) error {
	l := len(datas)
	keys := make([]interface{}, 0, l)
	values := make([]interface{}, 0, l)
	for i := 0; i < l; i ++ {
		data := datas[i]
		switch data.RewardType{
		case model.CoinType_Money:
			{
				keys = append(keys, "money")
				values = append(values, data.Num)
			}
		case model.CoinType_Gold:
			{
				keys = append(keys, "gold")
				values = append(values, data.Num)
			}
		}
	}
	keys = append(keys, values...)
	seelog.Debugf("UpdateCoins data:%v", datas)
	return accountProvider.UpdateAccountValue(gameid, openid, keys...)
}

func ChangeCoinTypeToDesc(addType int) string{
	switch addType{
	case model.ChangeCoinType_Advent:
		{
			return "观看广告"
		}
	case model.ChangeCoinType_FriendMatch:
		{
			return "邀请好友对战"
		}
	case model.ChangeCoinType_OneHundred:
		{
			return "百元赛"
		}
	case model.ChangeCoinType_RedPacket:
		{
			return "红包赛"
		}
	case model.ChangeCoinType_RandMatch:
		{
			return "随机匹配"
		}
	case model.ChangeCoinType_MatchTicket:
		{
			return "比赛门票"
		}
	case model.ChangeCoinType_InviteFriend:
		{
			return "邀请好友"
		}
	}
	return ""
}

func LogCoinsOnChange(gameid string, openid string, pregold uint64, premoney uint64, curgold uint64, curmoney uint64, desc string){
	infostr := fmt.Sprintf("用户[%s]", gameid + openid)
	if pregold <= curgold {
		infostr += fmt.Sprintf("获得金币:%d   ", curgold - pregold)
	}else{
		infostr += fmt.Sprintf("消耗金币:%d   ", pregold - curgold)
	}
	if premoney <= curmoney {
		infostr += fmt.Sprintf("获得现金:%d", curmoney - premoney)
	}else{
		infostr += fmt.Sprintf("消耗现金:%d", premoney - curmoney)
	}
	seelog.Debugf("%s,类型:%s", infostr, desc)
}

func UpdateQDInfo(chanid string, gameid string, num string) error{
	return accountProvider.UpdateQDInfo(gameid, chanid, num)
}