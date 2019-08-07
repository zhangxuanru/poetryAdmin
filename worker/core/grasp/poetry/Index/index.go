package Index

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
)

//抓取首页
type Index struct {
	content chan string
}

func NewIndex() *Index {
	return &Index{}
}

func (i *Index) GetAllData() {
	i.GetIndexHtml()
}

//首页-诗文分类
func (i *Index) GetPoetryCategory() {

}

//首页-名句分类
func (i *Index) GetPoetryFamousCategory() {

}

//首页-作者
func (i *Index) GetPoetryAuthor() {

}

//获取首页html内容
func (i *Index) GetIndexHtml() {
	var (
		resp  *http.Response
		err   error
		bytes []byte
	)
	if resp, bytes, err = tools.NewHttpReq().HttpGet(config.G_Conf.GuShiWenIndexUrl); err != nil {
		go data.G_GraspResult.PushError(err)
		return
	}

	logrus.Info("bytes:", string(bytes))
	logrus.Info("resp:", resp)
	return
}
