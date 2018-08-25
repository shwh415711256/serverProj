package main

import (
	"github.com/name5566/leaf"
	lconf "github.com/name5566/leaf/conf"
	"server/conf"
	"server/game"
	"server/gate"
	"server/login"
	"github.com/smallnest/rpcx/client"
)

var (
	dbxClient client.XClient
)

func main() {
	lconf.LogLevel = conf.Server.LogLevel
	lconf.LogPath = conf.Server.LogPath
	lconf.LogFlag = conf.LogFlag
	lconf.ConsolePort = conf.Server.ConsolePort
	lconf.ProfilePath = conf.Server.ProfilePath

	d := client.NewEtcdDiscovery("/zq/rpcx", conf.Server.DBServerPath, []string{conf.Server.EtcdAddr}, nil)
	dbxClient = client.NewXClient(conf.Server.DBServerPath, client.Failover, client.RandomSelect, d, client.DefaultOption)
	defer dbxClient.Close()

	leaf.Run(
		game.Module,
		gate.Module,
		login.Module,
	)
}
