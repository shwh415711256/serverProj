package main

import (
	"Resource/router"
	"Resource/conf"
	"os"
	"os/signal"
	"Resource/logger"
	"github.com/cihub/seelog"
	"Resource/router/middleware"
)

func main(){
	conf.Init()
	logger.Init(conf.ServerData.SeelogXmlPath)

	seelog.Debugf("start route file")
	fr := router.FileListenerLoad(middleware.LoggerMiddleware)
	go func() {
		err := fr.Run(conf.ServerData.FileListenPort)
		if err != nil {
			seelog.Errorf("FileServer start error[%v]", err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	sig := <-c
	seelog.Debugf("signal [%v]", sig)
	seelog.Flush()
}
