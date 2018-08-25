package version

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"GameServer/conf"
)

func TestReturnIgnore(c *gin.Context) {
	conf.VersionIgnore.Lock.RLock()
	defer conf.VersionIgnore.Lock.RUnlock()
	uri := c.Request.RequestURI
	if flag, find := conf.VersionIgnore.IgnoreData[uri]; find {
		c.Data(http.StatusOK, "text/plain", []byte(flag))
	}
}