package file

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"Resource/conf"
	"net/http"
	"github.com/cihub/seelog"
)

func HandlerUpLoad(c *gin.Context){
	form, err := c.MultipartForm()
	if err != nil {
		c.Data(500, "text/plain", []byte("form error" + err.Error()))
		return
	}
	dir := form.Value["dir"][0]
	dir = conf.ServerData.ResourcePath + "/" + dir + "/"

	isDirExitOrMkDir(dir)
	seelog.Debugf("%s", dir)
	files := form.File["files"]
	for _, file := range files {
		seelog.Debugf("upload file name:%s, header:%s", file.Filename, file.Header)
		fmt.Print(file.Filename)
		c.SaveUploadedFile(file, dir + file.Filename)
	}
	c.Data(http.StatusOK, "text/plain", []byte("ok"))
}

func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func HandlerDownLoad(c *gin.Context){

}
