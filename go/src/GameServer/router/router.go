package router

import (
	"github.com/gin-gonic/gin"
	"GameServer/module/account"
	"GameServer/module/version"
	"GameServer/conf"
	"GameServer/module/dbclient"
	"GameServer/module/room"
)

func GamelistenerLoad(middleware ...gin.HandlerFunc) *gin.Engine{
	r := gin.Default()
	r.Use(middleware...)

	// acount
	r.GET("/login", account.Login)
	r.POST("/updateWxUserInfo", account.UpdateWxUserInfo)
	r.GET("/getUserAccInfo", account.GetUserAccInfo)

	// version
	r.GET("/ignore", dbclient.GetVersionIgnoreInfo)
	r.GET("/ignorereload", dbclient.ReloadVersionIgnoreInfo)

	// room
	r.POST("/createRoom", room.CreateRoom)
	r.POST("/getOneMatch", room.RandomMatchOne)
	r.POST("/commitGameHis", room.CommitGameHis)

	// sign
	r.GET("/signInfo", )

	// match
	r.POST("/enrollMatch", room.EnrollMatch)

	// 渠道信息
	r.GET("/pages/index/index", account.UpdateQDLoginNum)

	// test 僵尸与消除 version_ignore 之后修改
	r.GET("/test", version.TestReturnIgnore)

	// test2 打枪
	r.GET("/test2", version.TestReturnIgnore)

	// g2048 2048
	r.GET("/g2048", version.TestReturnIgnore)

	// 临时的加载接口
	r.GET("/reloadignore", conf.ReloadVersion)

	return r
}