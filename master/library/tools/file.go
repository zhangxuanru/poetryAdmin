package tools

import (
	"os"
)

//判断文件是否存在
func PathExists(file string) (ret bool, err error) {
	if _, err = os.Stat(file); err != nil {
		return false, err
	}
	if os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}
