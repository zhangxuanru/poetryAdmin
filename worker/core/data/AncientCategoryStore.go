/*
@Time : 2019/9/2 15:37
@Author : zxr
@File : AncientCategoryStore
@Software: GoLand
*/
//保存古籍-栏目分类
package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/core/define"
	"time"
)

type AncientCategoryStore struct {
}

func NewAncientCategoryStore() *AncientCategoryStore {
	return &AncientCategoryStore{}
}

//载入古文首页分类数据
func (a *AncientCategoryStore) LoadCategoryData(data interface{}, params interface{}) {
	var (
		categoryList *[]define.GuWenCategoryList
		ok           bool
		err          error
		id           int64
	)
	if categoryList, ok = data.(*[]define.GuWenCategoryList); ok == false {
		logrus.Infoln("data conver GuWenCategoryList error")
		return
	}
	model := models.NewAncientCategory()
	for _, category := range *categoryList {
		insertData := &models.AncientCategory{
			CatName: category.CategoryName,
			SrcUrl:  category.LinkUrl,
			Pid:     0,
			Sort:    category.Sort,
			AddDate: time.Now().Unix(),
		}
		if id, err = model.SaveCategory(insertData); err != nil {
			logrus.Infoln("SaveCategory error:", err)
			continue
		}
		if id == 0 {
			logrus.Infoln("保存顶级分类失败......")
			continue
		}
		for _, nodeCat := range category.SubNode {
			insertData := &models.AncientCategory{
				CatName: nodeCat.CategoryName,
				SrcUrl:  nodeCat.LinkUrl,
				Sort:    nodeCat.Sort,
				Pid:     id,
				AddDate: time.Now().Unix(),
			}
			if _, err = model.SaveCategory(insertData); err != nil {
				logrus.Infoln("SaveCategory error:", err)
			}
		}
	}
}
