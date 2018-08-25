package router

import (
	"github.com/gin-gonic/gin"
	"Resource/module/file"
	"Resource/conf"
)

func FileListenerLoad(middleware ...gin.HandlerFunc) *gin.Engine{
	r := gin.Default()
	r.Use(middleware...)
	// HTML渲染
	r.LoadHTMLGlob("public/*")

	r.Static("/download", conf.ServerData.ResourcePath)

	// files
	fileGroup := r.Group("/files")
	fileGroup.GET("/upload", file.IndexHandler)
	fileGroup.POST("/upload", file.HandlerUpLoad)
	return r
}