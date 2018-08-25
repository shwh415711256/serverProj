package conf

import (
	"encoding/json"
	"github.com/name5566/leaf/log"
	"io/ioutil"
	"Common/model"

)

var Server struct {
	LogLevel    string
	LogPath     string
	WSAddr      string
	CertFile    string
	KeyFile     string
	TCPAddr     string
	MaxConnNum  int
	ConsolePort int
	ProfilePath string
	DBServerPath string
	EtcdAddr string
	RpcxPort string
	RoomServerPath string
}

var AllGameInfo struct {
	GameInfo []*OneGameInfo
}

type OneGameInfo struct {
	GameId string
	RoomInfo []*model.OneRoomInfo
}

func GetConfRoomInfo(gameId string, gameType int) *model.OneRoomInfo{
	gl := len(AllGameInfo.GameInfo)
	for i := 0; i < gl; i ++{
		if gameId == AllGameInfo.GameInfo[i].GameId {
			gameData := AllGameInfo.GameInfo[i]
			rl := len(gameData.RoomInfo)
			for j := 0; j < rl; j ++ {
				if gameType == gameData.RoomInfo[j].RoomType {
					return gameData.RoomInfo[j]
				}
			}
		}
	}
	return nil
}

func init() {
	//data, err := ioutil.ReadFile("F:/proj/go/src/RoomServer/src/server/conf/server.json")
	data, err := ioutil.ReadFile("conf/server.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(data, &Server)
	if err != nil {
		log.Fatal("%v", err)
	}
	//gamedata, err := ioutil.ReadFile("F:/proj/go/src/RoomServer/src/server/conf/room.json")
	gamedata, err := ioutil.ReadFile("conf/room.json")
	if err != nil {
		log.Fatal("%v", err)
	}
	err = json.Unmarshal(gamedata, &AllGameInfo)
	if err != nil {
		log.Fatal("%v", err)
	}
}
