package model

import (
	"sync"
)

const (
	GameId_JS			= "1"			// 僵尸  1
	GaemId_XC			= "2"			// 消除  2
	GameId_Gun			= "3"			// 打枪  3
	GameId_G2048		= "4"			// 2048  4
	GameId_baketBall 	= "5"			// 篮球 5
	GameId_Run			= "6"			// 跑圈 6
)

type RpcVersionInfoRequest struct {
	Key string
}

type RpcVersionInfoResponse struct {
	IgnoreFlag int
}

type VersionConfData struct {
	Lock sync.RWMutex
	IgnoreData map[string]int
}

type MatchConfData struct {
	Lock sync.RWMutex
	MatchData map[int]*OneMatchConfData
}

type OneMatchConfData struct{
	TicketNum int
	LoseNum int
	RewardData []MatchRewardInfo
}

type MatchRewardInfo struct {
	RankStart	int
	RankEnd		int
	Reward []OneReward
}

type OneReward struct {
	RewardType int
	Num int64
}

type SignRewardData struct{
	Lock sync.RWMutex
	RewardData map[int]*[][]OneReward
}

type CommonEmptyRequest struct {
}

type CommonEmptyResponse struct {
}