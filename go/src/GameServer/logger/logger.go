package logger

import "github.com/cihub/seelog"

func Init(filePath string){
	logger, err := seelog.LoggerFromConfigAsFile(filePath)
	if err != nil {
		seelog.Errorf("parse conf/seelog.xml error[%v]", err)
		return
	}
	seelog.ReplaceLogger(logger)
}
