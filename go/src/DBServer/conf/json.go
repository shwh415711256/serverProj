package conf

import (
	"io/ioutil"
	"encoding/json"
	"Common/model"
)

var ServerData struct {
	ListenPort string
	SeelogXmlPath string
	EtcdAddr string
	DBAddr string
	DBUserName string
	DBPassword string
	DBName string
	MaxOpenConns int
	MaxIdleConns int
	RedisMaxIdleNum int
	RedisMaxActiveNum int
	RedisIdleTimeOut int
	RedisAddr string
	RedisPass string
}

var ConfData struct {
	VersionData model.VersionConfData
	MatchData model.MatchConfData
	SignRewardData model.SignRewardData
}

func Init(){
	data, err := ioutil.ReadFile("conf/conf.online.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &ServerData)
	if err != nil {
		return
	}
	ConfInit()
}

func ConfInit(){
	ConfData.VersionData.IgnoreData = make(map[string]int)
	ConfData.MatchData.MatchData = make(map[int]*model.OneMatchConfData)
	ConfData.SignRewardData.RewardData = make(map[int]*[][]model.OneReward)
}
