package match

import (
	"DBServer/store/match"
	"Common/model"
	"DBServer/conf"
	"github.com/cihub/seelog"
	"strconv"
	"strings"
)

var (
	matchProvider match.MatchProvider
)

func Init(){
	matchProvider = &match.MatchDefaultProvider{}
	matchProvider.(*match.MatchDefaultProvider).Init()
	LoadMatchConfig(&conf.ConfData.MatchData)
}

func LoadMatchConfig(data *model.MatchConfData){
	rows, err := matchProvider.LoadMatchConfigData()
	if err != nil {
		seelog.Errorf("LoadMatchConfig error[%v]", err)
		return
	}
	for _, r := range *rows {
		oneMatch := &model.OneMatchConfData{}
		matchType, err := strconv.ParseInt(string(r["match_type"]), 10, 32)
		if err != nil {
			seelog.Errorf("LoadMatchConfig parse matchType error:%v", err)
			continue
		}
		ticketNum, err := strconv.ParseInt(string(r["ticket"]), 10, 32)
		if err != nil {
			seelog.Errorf("LoadMatchConfig parse ticketNum error:%v", err)
			continue
		}
		oneMatch.TicketNum = int(ticketNum)
		loseNum, err := strconv.ParseInt(string(r["losenum"]), 10, 32)
		if err != nil {
			seelog.Errorf("LoadMatchConfig parse loseNum error:%v", err)
			continue
		}
		oneMatch.LoseNum = int(loseNum)
		rewardInfos := strings.Split(string(r["reward"]),";")
		l := len(rewardInfos)
		for i := 0; i < l; i ++ {
			reward := &model.MatchRewardInfo{}
			strs := strings.Split(rewardInfos[i], "|")
			if len(strs) != 3 {
				seelog.Errorf("LoadMatchConfig parse reward error:%v", err)
				continue
			}
			ranks, err := strconv.ParseInt(strs[0], 10, 32)
			if err != nil {
				seelog.Errorf("LoadMatchConfig parse rankstart error:%v", err)
				continue
			}
			ranke, err := strconv.ParseInt(strs[1], 10, 32)
			if err != nil {
				seelog.Errorf("LoadMatchConfig parse rankend error:%v", err)
				continue
			}
			reward.RankStart = int(ranks)
			reward.RankEnd = int(ranke)
			rewards := strings.Split(strs[2], ",")
			for j := 0; j < len(rewards); j ++ {
				info := strings.Split(rewards[j], "-")
				if len(info) != 2 {
					seelog.Errorf("LoadMatchConfig parse rewards error:%v", err)
					continue
				}
				rtype, err := strconv.ParseInt(info[0], 10, 32)
				if err != nil {
					seelog.Errorf("LoadMatchConfig parse rewardType error:%v", err)
					continue
				}
				num, err := strconv.ParseInt(info[1], 10, 32)
				if err != nil {
					seelog.Errorf("LoadMatchConfig parse rewardNum error:%v", err)
					continue
				}
				reward.Reward = append(reward.Reward, model.OneReward{
					RewardType: int(rtype),
					Num:        num,
				})
			}
			oneMatch.RewardData = append(oneMatch.RewardData, *reward)
		}
		data.MatchData[int(matchType)] = oneMatch
	}
}

func ReloadMatchConfig(data *model.MatchConfData){
	data.Lock.Lock()
	defer data.Lock.Unlock()
	LoadMatchConfig(data)
}