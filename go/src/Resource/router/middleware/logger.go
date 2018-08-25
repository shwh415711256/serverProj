package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/cihub/seelog"
	"time"
)

func LoggerMiddleware(c *gin.Context){
	seelog.Debugf("time:%s, url:%s, method:%s", time.Now().String(), c.Request.URL, c.Request.Method)
}
