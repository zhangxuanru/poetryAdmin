package Index

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/grasp/poetry/base"
)

//抓取首页
type Index struct {
	Content string
	GoQuery *goquery.Document
}

func NewIndex() *Index {
	return &Index{}
}

//获取首页所有内容
func (i *Index) GetAllData() {
	logrus.Info("GetAllData start .......")
	if err := i.GetIndexSource(); err != nil {
		logrus.Debug("GetIndexHtml err:", err)
		return
	}
	if base.CheckContent(i.Content) == false {
		logrus.Debug("CheckContent err: content is nil")
		return
	}
	i.GetPoetryCategory()
	i.GetPoetryFamousCategory()
	i.GetPoetryAuthor()

	//通过goquery 解析 html
	i.GoQuery.Index()
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
func (i *Index) GetIndexSource() (err error) {
	var (
		query *goquery.Document
		bytes []byte
	)
	if bytes, err = base.GetHtml(config.G_Conf.GuShiWenIndexUrl); err != nil {
		return
	}
	if len(bytes) > 0 {
		i.Content = string(bytes)
		query, err = tools.NewDocumentFromReader(i.Content)
	}
	if err != nil {
		return err
	}
	i.GoQuery = query
	return nil
}
