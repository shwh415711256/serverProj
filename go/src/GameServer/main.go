package main

import (
	"github.com/cihub/seelog"
	"os"
	"os/signal"
	"GameServer/router"
	"GameServer/conf"
	"GameServer/logger"
	"GameServer/router/middleware"
	"GameServer/module/dbclient"
	"GameServer/module/room"
)

func main(){
	conf.Init()
	logger.Init(conf.ServerData.SeelogXmlPath)

	seelog.Debugf("start route file")
	gr := router.GamelistenerLoad(middleware.LoggerMiddleware)
	go func() {
		err := gr.Run(conf.ServerData.ListenPort)
		if err != nil {
			seelog.Errorf("GameServer start error[%v]", err)
		}
	}()

	dbclient.Init()
	defer dbclient.GetDBManager().Close()

	room.Init()
	defer room.GetRoomManager().Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	seelog.Debugf("signal [%v]", sig)
	seelog.Flush()
}
