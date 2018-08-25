package conf

import (
	"io/ioutil"
	"encoding/json"
	"sync"
	"github.com/cihub/seelog"
	"github.com/gin-gonic/gin"
	"net/http"
)

var ServerData struct {
	ListenPort string
	SeelogXmlPath string
	EtcdAddr string
	DBServerPath string
	RoomServerPath string
}

var VersionIgnore struct {
	Lock sync.RWMutex
	IgnoreData map[string]string
}

func ReloadVersion(c *gin.Context) {
	VersionIgnore.Lock.Lock()
	defer VersionIgnore.Lock.Unlock()
	idata, err := ioutil.ReadFile("conf/version.ignore.json")
	if err != nil {
		seelog.Errorf("ReadFile conf.ignore.json error[%v]", err)
		return
	}
	err = json.Unmarshal(idata, &VersionIgnore.IgnoreData)
	if err != nil {
		seelog.Errorf("VersionIgnore unmarshal error[%v]", err)
		return
	}
	c.Data(http.StatusOK, "text/plain", []byte("ok"))
}

func Init(){
	data, err := ioutil.ReadFile("conf/conf.online.json")
	if err != nil {
		seelog.Errorf("ReadFile conf.online.json error[%v]", err)
		return
	}
	err = json.Unmarshal(data, &ServerData)
	if err != nil {
		seelog.Errorf("ServerData unmarshal error[%v]", err)
		return
	}
	idata, err := ioutil.ReadFile("conf/version.ignore.json")
	if err != nil {
		seelog.Errorf("ReadFile conf.ignore.json error[%v]", err)
		return
	}
	err = json.Unmarshal(idata, &VersionIgnore.IgnoreData)
	if err != nil {
		seelog.Errorf("VersionIgnore unmarshal error[%v]", err)
		return
	}
}
