package game

import (
	"DBServer/store/game"
	"errors"
	"math/rand"
)

var (
	gameProvider game.GameProvider
)

func Init() {
	gameProvider = &game.GameDefaultProvider{}
	gameProvider.(*game.GameDefaultProvider).Init()
}

func RandOneMatchInfo(gameid string, openid string) (string, string, error){
	rows, err := gameProvider.LoadOneMatchHis(gameid, openid)
	if err != nil {
		return "", "", err
	}
	if len(*rows) == 0{
		return "", "", errors.New("db have no data")
	}
	index := rand.Intn(len(*rows))
	r := (*rows)[index]
	return string(r["openid"]), string(r["hisdata"]), nil
}

func InsertGameHis(gameid string, openid string, data string) error {
	return gameProvider.InsertOneMatchHis(gameid, openid, data)
}