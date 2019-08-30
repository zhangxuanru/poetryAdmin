package base

import (
	"errors"
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
)

//获取html内容
func GetHtml(url string) (bytes []byte, err error) {
	//logrus.Infoln("GetHtml.....", url)
	if _, bytes, err = tools.NewHttpReq().HttpGet(url); err != nil {
		go data.G_GraspResult.PushError(err)
		return
	}
	return bytes, err
}

//检查content是否有内容
func CheckContent(content string) (ret bool) {
	if len(content) == 0 {
		return false
	}
	return true
}

//读取测试文件内容
func GetTestFile(file string) (bytes []byte, err error) {
	logrus.Infoln("GetTestFile.....")
	if ret, _ := tools.PathExists(file); ret == true {
		return tools.ReadFile(file)
	}
	return nil, errors.New(file + "file is not exists")
}
