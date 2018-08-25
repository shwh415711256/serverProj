package rpcx

import (
	"context"
	"Common/model"
	"DBServer/conf"
	"DBServer/module/version"
	"DBServer/db/cache"
	"DBServer/module/account"
	"github.com/cihub/seelog"
	"DBServer/module/game"
	"errors"
	"fmt"
)

type DBProcess struct{}

func (p *DBProcess) GetVersionInfo(ctx context.Context, args *model.RpcVersionInfoRequest, reply *model.RpcVersionInfoResponse) error {
	conf.ConfData.VersionData.Lock.RLock()
	defer conf.ConfData.VersionData.Lock.RUnlock()
	key := args.Key
	ignoreFlag, find := conf.ConfData.VersionData.IgnoreData[key]
	if find {
		reply.IgnoreFlag = ignoreFlag
	}else {
		reply.IgnoreFlag = 0
	}
	return nil
}

func (p *DBProcess) ReloadVersionInfo(ctx context.Context, args *model.CommonEmptyRequest, reply *model.CommonEmptyResponse) error {
	version.ReloadVersionConfig(&conf.ConfData.VersionData)
	return nil
}

func (p *DBProcess) RpcLoadAccountInfo(ctx context.Context, args *model.RpcLoadAccountInfoRequest, reply *model.RpcLoadAccountInfoResponse) error{
	info, find := cache.GetAccountInfo(args.GameId, args.OpenId)
	if !find {
		info = account.NewAccount(args.GameId, args.OpenId)
		/* 正式
		if args.InviteOpenid != "" {
			inviteInfo, ifind := cache.GetAccountInfo(args.GameId, args.InviteOpenid)
			if ifind{
				pregold := inviteInfo.Gold
				inviteInfo.Gold += 5
				account.UpdateCoins(args.GameId, args.InviteOpenid, []model.OneReward{
					model.OneReward{RewardType:model.CoinType_Gold,Num:int64(inviteInfo.Gold)},
					model.OneReward{RewardType:model.CoinType_Money,Num:int64(inviteInfo.Money)},
				})
				cache.DelAccountInfo(args.GameId, args.InviteOpenid)
				account.LogCoinsOnChange(args.GameId, args.InviteOpenid, pregold, inviteInfo.Money, inviteInfo.Gold, inviteInfo.Money,
					account.ChangeCoinTypeToDesc(model.ChangeCoinType_InviteFriend))
			}
		}
		*/
		cache.SetAccountInfo(args.GameId, info)
	}
	// 测试用
	if args.InviteOpenid != "" {
		inviteInfo, ifind := cache.GetAccountInfo(args.GameId, args.InviteOpenid)
		if ifind{
			pregold := inviteInfo.Gold
			inviteInfo.Gold += 5
			account.UpdateCoins(args.GameId, args.InviteOpenid, []model.OneReward{
				model.OneReward{RewardType:model.CoinType_Gold,Num:int64(inviteInfo.Gold)},
				model.OneReward{RewardType:model.CoinType_Money,Num:int64(inviteInfo.Money)},
			})
			cache.DelAccountInfo(args.GameId, args.InviteOpenid)
			account.LogCoinsOnChange(args.GameId, args.InviteOpenid, pregold, inviteInfo.Money, inviteInfo.Gold, inviteInfo.Money,
				account.ChangeCoinTypeToDesc(model.ChangeCoinType_InviteFriend))
		}
	}
	reply.AccountInfoData = *info
	if (args.NeedWxData == "1") {
		wxinfo, err := account.LoadWxUserDataFromDB(args.GameId, args.OpenId)
		if err != nil {
			seelog.Errorf("loadwxinfo from db error[%v]", err)
			return nil
		}
		reply.UserData = *wxinfo
	}
	return nil
}

func (p *DBProcess) RpcUpdateWxUserInfo(ctx context.Context, args *model.RpcUpdateWxInfoRequest, reply *model.RpcCommonEmptyResponse) error{
	err := account.UpdateWxInfo(args.GameId, &args.WxData)
	if err != nil {
		seelog.Errorf("UpdateWxUserInfo error[%v]",err)
		return err
	}
	return nil
}

func (p *DBProcess) RpcGetOneMatchInfo(ctx context.Context, args *model.RpcGetOneMatchInfoRequest, reply *model.RpcGetOneMatchInfoResponse) error {
	openid, data , err := game.RandOneMatchInfo(args.Gameid, args.Openid)
	if err != nil {
		seelog.Errorf("RpcGetOneMatchInfo error[%v], gameid:%s, openid:%s", err, args.Gameid, args.Openid)
		return err
	}
	reply.Data = data
	accInfo, find := cache.GetAccountInfo(args.Gameid, openid)
	if !find {

		err = errors.New(fmt.Sprintf("have no accountInfo, gameid:%s, openid:%s", args.Gameid, openid))
		seelog.Errorf("RpcGetOneMatchInfo error[%v]", err)
		return err
	}
	wxinfo, err := account.LoadWxUserDataFromDB(args.Gameid, openid)
	if err != nil {
		seelog.Errorf("loadwxinfo from db error[%v]", err)
		return nil
	}
	reply.AccData = *accInfo
	reply.WxData = *wxinfo
	return nil
}

func (p *DBProcess) RpcInsertGameHis(ctx context.Context, args *model.GameHisInfo, reply *model.RpcCommonEmptyResponse) error {
	err := game.InsertGameHis(args.Gameid, args.Openid, args.Data)
	if err != nil {
		seelog.Errorf("RpcInsertGameHis error[%v]", err)
		return err
	}
	return nil
}

func (p *DBProcess) RpcDelMatchTicket(ctx context.Context, args *model.RpcDelMatchTicketRequest, reply *model.RpcDelMatchTicketResponse) error {
	conf.ConfData.MatchData.Lock.RLock()
	defer conf.ConfData.MatchData.Lock.RUnlock()
	matchData, find := conf.ConfData.MatchData.MatchData[args.MatchType]
	if !find {
		reply.Result = model.MatchOptRet_NoMatchData
		return errors.New(fmt.Sprintf("RpcDelMatchTicket not find matchData, matchType:%d", args.MatchType))
	}
	info, find := cache.GetAccountInfo(args.Gameid, args.Openid)
	if !find {
		reply.Result = model.MatchOptRet_NoPlayer
		return errors.New(fmt.Sprintf("player not exist,gameid:%s,openid:%s", args.Gameid, args.Openid))
	}
	if info.Gold < uint64(matchData.TicketNum+matchData.LoseNum) {
		reply.Result = model.MatchOptRet_GoldNotEnough
		return errors.New(fmt.Sprintf("player gold not enough, gold:%d, neednum:%d", info.Gold, matchData.LoseNum + matchData.TicketNum))
	}
	pregold := info.Gold
	premoney := info.Money
	info.Gold -= uint64(matchData.TicketNum)
	err := account.UpdateCoins(args.Gameid, args.Openid, []model.OneReward{
		model.OneReward{RewardType:model.CoinType_Gold, Num:int64(info.Gold)},
	})
	if err != nil {
		reply.Result = model.MatchOptRet_DelError
		return err
	}
	cache.DelAccountInfo(args.Gameid, args.Openid)
	account.LogCoinsOnChange(args.Gameid, args.Openid, pregold, premoney, info.Gold, info.Money,
		account.ChangeCoinTypeToDesc(model.ChangeCoinType_MatchTicket))
	reply.Result = model.MatchOptRet_Success
	reply.CurNum = info.Gold
	return nil
}

func (p *DBProcess) RpcCalMatchReward(ctx context.Context, args *model.RpcCalMatchRewardRequest, reply *model.RpcCalMatchRewardResponse) error{
	conf.ConfData.MatchData.Lock.RLock()
	defer conf.ConfData.MatchData.Lock.RUnlock()
	matchData, find := conf.ConfData.MatchData.MatchData[args.MatchType]
	if !find {
		reply.Result = model.MatchOptRet_NoMatchData
		return errors.New(fmt.Sprintf("RpcCalMatchReward not find matchData, matchType:%d", args.MatchType))
	}
	info, find := cache.GetAccountInfo(args.Gameid, args.Openid)
	if !find {
		reply.Result = model.MatchOptRet_NoPlayer
		return errors.New(fmt.Sprintf("RpcCalMatchReward player not exist,gameid:%s,openid:%s", args.Gameid, args.Openid))
	}
	l := len(matchData.RewardData)
	pregold := info.Gold
	premoney := info.Money
	for i := 0; i < l; i ++{
		rdata := matchData.RewardData[i]
		if args.Rank >= rdata.RankStart && args.Rank <= rdata.RankEnd {
			rl := len(matchData.RewardData[i].Reward)
			for j := 0; j < rl; j ++{
				reward := matchData.RewardData[i].Reward[j]
				account.CalculateCoins(info, &reward)
			}
			break
		}
	}
	if args.Rank != 1 {
		if info.Gold > uint64(matchData.LoseNum) {
			info.Gold -= uint64(matchData.LoseNum)
		}else{
			seelog.Debugf("error gold not enough, gold:%d, lose:%d", info.Gold, matchData.LoseNum)
			info.Gold = 0
		}
	}
	err := account.UpdateCoins(args.Gameid, args.Openid, []model.OneReward{
		model.OneReward{RewardType:model.CoinType_Gold,Num:int64(info.Gold)},
		model.OneReward{RewardType:model.CoinType_Money,Num:int64(info.Money)},
	})
	if err != nil {
		reply.Result = model.MatchOptRet_DelError
		return err
	}
	cache.DelAccountInfo(args.Gameid, args.Openid)
	reply.Result = model.MatchOptRet_Success
	reply.CurGold = info.Gold
	reply.CurMoney = info.Money
	changeCoinInfo := account.ChangeCoinTypeToDesc(args.MatchType)
	changeCoinInfo += fmt.Sprintf(",名次:%d", args.Rank)
	account.LogCoinsOnChange(args.Gameid, args.Openid, pregold, premoney, info.Gold, info.Money, changeCoinInfo)
	return nil
}

func (p *DBProcess) RpcAddCoins(ctx context.Context, args *model.RpcAddCoinsRequest, reply *model.RpcAddCoinsResponse) error{
	info, find := cache.GetAccountInfo(args.Gameid, args.Openid)
	if !find {
		reply.Result = 0
		return errors.New(fmt.Sprintf("RpcAddCoins player not exist,gameid:%s,openid:%s", args.Gameid, args.Openid))
	}
	pregold := info.Gold
	premoney := info.Money
	l := len(args.Data)
	for i := 0; i < l; i ++ {
		account.CalculateCoins(info, &args.Data[i])
	}
	seelog.Debugf("pregold:%d, premoney:%d, curgold:%d, curmoney:%d", pregold, premoney, info.Gold, info.Money)
	err := account.UpdateCoins(args.Gameid, args.Openid, []model.OneReward{
		model.OneReward{RewardType:model.CoinType_Gold,Num:int64(info.Gold)},
		model.OneReward{RewardType:model.CoinType_Money,Num:int64(info.Money)},
	})
	if err != nil {
		reply.Result = 0
		return err
	}
	cache.DelAccountInfo(args.Gameid, args.Openid)
	reply.Result = 1
	account.LogCoinsOnChange(args.Gameid, args.Openid, pregold, premoney, info.Gold, info.Money, account.ChangeCoinTypeToDesc(model.ChangeCoinType_FriendMatch))
	return nil
}

func (p *DBProcess) RpcUpdateQDInfo(ctx context.Context, args *model.RpcUpdateQDInfoRequest, reply *model.RpcCommonEmptyResponse) error{
	return account.UpdateQDInfo(args.ChanId, args.Gameid, args.Num)
}

func (p *DBProcess) RpcGetSignInfo(ctx context.Context, args *model.RpcGetSignInfoRequest, reply *model.RpcGetSignInfoResponse) error{
	find := cache.GetSignInfo(args.GameId, args.OpenId, args.SignType, reply)
	if !find {
		return errors.New("")
	}
	return nil
}