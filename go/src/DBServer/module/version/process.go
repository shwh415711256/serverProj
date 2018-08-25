package version

import (
	"Common/model"
	"DBServer/store/version"
	"github.com/cihub/seelog"
	"strconv"
	"DBServer/conf"
)

var (
	versionProvider version.VersionProvider
)

func Init() {
	versionProvider = &version.VersionDefaultProvider{
	}
	versionProvider.(*version.VersionDefaultProvider).Init()
	LoadVersionConfig(&conf.ConfData.VersionData)
}

func LoadVersionConfig(data *model.VersionConfData) {
	rows, err := versionProvider.LoadVersionConfigData()
	if err != nil {
		seelog.Errorf("LoadVersionConfig error[%v]", err)
		return
	}
	for _, row := range *rows {
		key := row["key"]
		ignoreStr := row["ignore_flag"]
		ignore, err := strconv.ParseInt(string(ignoreStr), 10, 32)
		if err != nil {
			seelog.Errorf("LoadVersionConfig parse ignore_flag error[%v]", err)
			continue
		}
		data.IgnoreData[string(key)] = int(ignore)
	}
}

func ReloadVersionConfig(data *model.VersionConfData) {
	data.Lock.Lock()
	defer data.Lock.Unlock()
	data.IgnoreData = make(map[string]int)
	LoadVersionConfig(data)
}