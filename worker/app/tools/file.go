package tools

import (
	"io/ioutil"
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

//读取文件内容
func ReadFile(file string) (bytes []byte, err error) {
	return ioutil.ReadFile(file)
}
