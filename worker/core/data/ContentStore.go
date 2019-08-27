package data

import "github.com/sirupsen/logrus"

//保存诗词正文和赏析注释结果...
type contentStore struct {
}

func NewContentStore() *contentStore {
	return new(contentStore)
}

//保存诗词详情信息与资料
func (c *contentStore) LoadPoetryContentData(data interface{}, params interface{}) {
	/*
	 1.写入poetry_detail_category 诗词分类表，
	 2.写入poetry_detail_notes 文本内容表
	 3.写入poetry_content_trans 翻译表
	 4.写入poetry_content_apprec赏析表
	 5.更新poetry_content_relation表
	 6.更新poetry_content诗词内容
	*/
	logrus.Infoln("LoadPoetryContentData.............")
	logrus.Infof("%+v\n", data)
	logrus.Infof("%+v\n", params)
}
