package sign

import (
	"DBServer/store/sign"
	"Common/model"
	"DBServer/conf"
	"github.com/cihub/seelog"
	"strconv"
	"strings"
	"errors"
)

var (
	signProvider sign.SignProvider
)

func Init(){
	signProvider = &sign.SignDefaultProvider{}
	signProvider.(*sign.SignDefaultProvider).Init()
	LoadSignRewardConfig(&conf.ConfData.SignRewardData)
}

func LoadSignRewardConfig(data *model.SignRewardData){
	rows, err := signProvider.LoadSignRewardConfig()
	if err != nil {
		seelog.Errorf("LoadSignRewardConfig error[%v]", err)
		return
	}
	for _, r := range *rows {
		rewardData := [][]model.OneReward{}
		sign_type, err := strconv.ParseInt(string(r["type"]), 10, 32)
		if err != nil {
			seelog.Errorf("LoadSignRewardConfig  parse type error[%v]", err)
			continue
		}
		rewardStr := strings.Split(string(r["reward"]),";")
		rl := len(rewardStr)
		for i:= 0; i < rl; i ++{
			dayRewards := strings.Split(rewardStr[i], "|")
			if len(dayRewards) != 2 {
				seelog.Errorf("LoadSignRewardConfig parse reward error[len not match]")
				return
			}
			day, err := strconv.ParseInt(dayRewards[0], 10, 32)
			if err != nil {
				seelog.Errorf("LoadSignRewardConfig  parse day error[%v]", err)
				continue
			}
			rewards := strings.Split(dayRewards[1], ",")
			l := len(rewards)
			rewardData[day] = make([]model.OneReward,0,l)
			for j := 0; j < l; j ++{
				oneReward := strings.Split(rewards[j], "-")
				if len(oneReward) != 2 {
					seelog.Errorf("LoadSignRewardConfig  parse oneReward error[len not match]")
					continue
				}
				rewardType, err := strconv.ParseInt(oneReward[0], 10, 32)
				if err != nil {
					seelog.Errorf("LoadSignRewardConfig  parse rewardType error[%v]", err)
					continue
				}
				rewardNum, err := strconv.ParseInt(oneReward[1], 10, 32)
				if err != nil {
					seelog.Errorf("LoadSignRewardConfig  parse rewardNum error[%v]", err)
					continue
				}
				rewardData[day] = append(rewardData[day], model.OneReward{
					RewardType:int(rewardType),
					Num:rewardNum,
				})
			}
		}
		data.RewardData[int(sign_type)] = &rewardData
	}
}

func ReloadSignRewardConfig(){
	conf.ConfData.SignRewardData.Lock.Lock()
	defer conf.ConfData.SignRewardData.Lock.Unlock()
	LoadSignRewardConfig(&conf.ConfData.SignRewardData)
}

func LoadSignInfoFromDB(gameId string, openId string, signType int) (*model.RpcGetSignInfoResponse, error) {
	rows, err := signProvider.LoadSignInfo(gameId, openId, signType)
	if err != nil {
		return nil, err
	}
	if len(*rows) == 0 {
		return nil, errors.New("LoadSignInfoFromDB error[no data]")
	}
	resp := &model.RpcGetSignInfoResponse{}
	r := (*rows)[0]
	signSumDay, err := strconv.ParseInt(string(r["sign_day"]), 10, 32)
	if err != nil {
		return nil, err
	}
	isTodaySign := string(r["today_sign"])

	resp.SignSumDay = int(signSumDay)
	if isTodaySign == "1" {
		resp.IsTodaySign = true
	} else {
		resp.IsTodaySign = false
	}
	return resp, nil
}