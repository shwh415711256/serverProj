package file

import (
	"strings"
	"os"
)

func isDirExitOrMkDir(path string) {
	strs := strings.Split(path, "/")
	filePath := strs[0]
	l := len(strs)
	for i := 1; i < l; i ++ {
		filePath += "/" + strs[i]
		_, err := os.Stat(filePath)
		if err != nil && os.IsNotExist(err) {
			os.Mkdir(filePath, os.ModePerm)
		}
	}
}