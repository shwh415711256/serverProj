package rpcclient

import (
	"github.com/smallnest/rpcx/client"
	"server/conf"
	"github.com/name5566/leaf/log"
)

var(
	DbXclient client.XClient
)

func Init() {
	d := client.NewEtcdDiscovery("/zq/rpcx", conf.Server.DBServerPath, []string{conf.Server.EtcdAddr}, nil)
	DbXclient = client.NewXClient(conf.Server.DBServerPath, client.Failover, client.RandomSelect, d, client.DefaultOption)
	log.Debug("DbXclient connect success")
}

func Close() {
	DbXclient.Close()
}
