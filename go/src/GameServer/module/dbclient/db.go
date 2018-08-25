package dbclient

import (
	"github.com/smallnest/rpcx/client"
	"GameServer/conf"
)

type DBManager struct{
	dbxclient client.XClient
}

var (
	dbManager *DBManager
)

func GetDBManager() *DBManager {
	return dbManager
}

func Init(){
	d := client.NewEtcdDiscovery("/zq/rpcx", conf.ServerData.DBServerPath, []string{conf.ServerData.EtcdAddr}, nil)
	dbManager = &DBManager{
		dbxclient:client.NewXClient(conf.ServerData.DBServerPath, client.Failover, client.RandomSelect, d, client.DefaultOption),
	}
}

func (m *DBManager) Close() {
	m.dbxclient.Close()
}

func (m *DBManager) GetDbxclient() client.XClient{
	return m.dbxclient
}