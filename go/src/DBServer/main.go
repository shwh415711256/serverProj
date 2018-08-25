package main

import (
	"DBServer/conf"
	"os"
	"os/signal"
	"DBServer/logger"
	"github.com/cihub/seelog"
	"github.com/smallnest/rpcx/server"
	"github.com/smallnest/rpcx/serverplugin"
	"github.com/rcrowley/go-metrics"
	"time"
	"log"
	"DBServer/db/rpcx"
	"DBServer/module/version"
	"DBServer/db"
	"DBServer/module/account"
	"DBServer/module/game"
	"math/rand"
	"DBServer/module/match"
)

func main(){
	conf.Init()
	db.Init()
	account.Init()
	version.Init()
	match.Init()
	game.Init()
	rand.Seed(time.Now().UnixNano())
	logger.Init(conf.ServerData.SeelogXmlPath)

	s := server.NewServer()
	addRegistryPlugin(s)

	err := s.RegisterName("DBServer", new(rpcx.DBProcess), "")
	if err != nil {
		seelog.Errorf("%v", err)
	}
	go func() {
		err := s.Serve("tcp", conf.ServerData.ListenPort)
		if err != nil {
			seelog.Errorf("%v", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	seelog.Debugf("signal [%v]", sig)
	seelog.Flush()
}

// 注册到服务发现
func addRegistryPlugin(server *server.Server) {
	etcd := &serverplugin.EtcdRegisterPlugin{
		ServiceAddress:"tcp@127.0.0.1" + conf.ServerData.ListenPort,
		EtcdServers:[]string{conf.ServerData.EtcdAddr},
		BasePath:"/zq/rpcx",
		Metrics:metrics.NewRegistry(),
		UpdateInterval:time.Minute,
	}
	err := etcd.Start()
	if err != nil {
		log.Fatal(err)
	}
	server.Plugins.Add(etcd)
}