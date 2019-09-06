/*
@Time : 2019/9/6 19:15
@Author : zxr
@File : FamousStore
@Software: GoLand
*/
package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"time"
)

//保存名句结果...
type FamousStorage struct {
}

func NewFamousStorage() *FamousStorage {
	return new(FamousStorage)
}

//载入名句分类数据
func (f *FamousStorage) LoadFamousClassifyData(data interface{}, params interface{}) {
	var (
		classifyList *[]define.Classify
		categoryData models.Category
		ok           bool
		err          error
	)
	if classifyList, ok = data.(*[]define.Classify); ok == false {
		logrus.Infoln("LoadFamousClassifyData data conver to Classify error")
		return
	}
	for _, classify := range *classifyList {
		if categoryData, err = models.GetDataByCateName(classify.Title, 2); err != nil {
			logrus.Infoln("GetDataByCateName error:", err)
			continue
		}
		if categoryData.Id > 0 {
			logrus.Infoln(classify.Title, "已存在")
			continue
		}
		saveData := &models.Category{
			CatName:        classify.Title,
			SourceUrl:      classify.LinkUrl,
			SourceUrlCrc32: tools.Crc32(classify.LinkUrl),
			ShowPosition:   2,
			AddDate:        time.Now().Unix(),
		}
		if _, err = models.SaveCategoryData(saveData); err != nil {
			logrus.Infoln("SaveCategoryData error:", err)
			continue
		}
	}
	logrus.Infoln("LoadFamousClassifyData save ok")
}
