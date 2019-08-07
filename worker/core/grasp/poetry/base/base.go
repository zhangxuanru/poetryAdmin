package base

import (
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
)

//获取html内容
func GetHtml(url string) (bytes []byte, err error) {
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
