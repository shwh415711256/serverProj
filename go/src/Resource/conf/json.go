package conf

import (
	"io/ioutil"
	"encoding/json"
)

var ServerData struct {
	FileListenPort string
	ResourcePath string
	SeelogXmlPath string
}

func Init(){
	data, err := ioutil.ReadFile("conf/conf.online.json")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &ServerData)
	if err != nil {
		return
	}
}
