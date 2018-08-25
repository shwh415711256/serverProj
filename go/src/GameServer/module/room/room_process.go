package room

import (
	"github.com/smallnest/rpcx/client"
	"GameServer/conf"
	"Common/model"
	"sync"
)

type RoomManager struct{
	roomXclient client.XClient
	randMatchLock sync.RWMutex
	randMatchInfoMap map[string]*model.RandMatchEnemy
}

var (
	roomManager *RoomManager
)

func GetRoomManager() *RoomManager {
	return roomManager
}

func Init(){
	d := client.NewEtcdDiscovery("/zq/rpcx", conf.ServerData.RoomServerPath, []string{conf.ServerData.EtcdAddr}, nil)
	roomManager = &RoomManager{
		roomXclient:client.NewXClient(conf.ServerData.RoomServerPath, client.Failover, client.RandomSelect, d, client.DefaultOption),
		randMatchInfoMap:make(map[string]*model.RandMatchEnemy),
	}
}

func (m *RoomManager) Close() {
	m.roomXclient.Close()
}

func (m *RoomManager) AddRandMatchEnemy(gameid string, openid string, data *model.RandMatchEnemy){
	m.randMatchLock.Lock()
	defer m.randMatchLock.Unlock()
	m.randMatchInfoMap[gameid+openid] = data
}

func (m *RoomManager) RemoveRandMatchEnemy(gameid string, openid string) {
	m.randMatchLock.Lock()
	defer m.randMatchLock.Unlock()
	delete(m.randMatchInfoMap, gameid+openid)
}

func (m *RoomManager) GetRandMatchEnemy(gameid string, openid string) (model.RandMatchEnemy, bool){
	m.randMatchLock.RLock()
	defer m.randMatchLock.RUnlock()
	data, find := m.randMatchInfoMap[gameid+openid]
	if !find {
		return model.RandMatchEnemy{}, false
	}
	return *data, true
}
