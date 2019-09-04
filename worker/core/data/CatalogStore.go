/*
@Time : 2019/9/4 18:18
@Author : zxr
@File : CatalogStore
@Software: GoLand
*/
package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/core/define"
	"time"
)

//保存古籍-书籍目录信息
type BookCatalogueStore struct {
}

func NewBookCatalogueStore() *BookCatalogueStore {
	return &BookCatalogueStore{}
}

//载入目录数据
func (b *BookCatalogueStore) LoadCatalogueData(data interface{}, params interface{}) {
	var (
		catId           int64
		err             error
		cataLogData     []*define.CataLogData
		ok              bool
		bookCover       *define.GuWenBookCover
		cataCategoryObj *models.BookCatalogueCategory
		catalogueObj    *models.BookCatalogue
		categoryData    models.AncientCategory
		bookData        models.AncientBook
	)
	if cataLogData, ok = data.([]*define.CataLogData); ok == false {
		logrus.Infoln("LoadCatalogueData data conver to CataLogData error")
		return
	}
	if bookCover, ok = params.(*define.GuWenBookCover); ok == false {
		logrus.Infoln("LoadCatalogueData params conver to GuWenBookCover error")
		return
	}
	categoryData, err = models.NewAncientCategory().GetCategoryDataByName(bookCover.GuWenCategory.CategoryName)
	if err != nil || categoryData.Id == 0 {
		logrus.Infoln("GetCategoryDataByName err:", err)
		return
	}
	if bookData, err = models.NewAncientBook().GetBookByTitleAndCatId(bookCover.Title, categoryData.Id); err != nil {
		logrus.Infoln("GetBookByTitleAndCatId error:", err)
		return
	}
	cataCategoryObj = models.NewBookCatalogueCategory()
	catalogueObj = models.NewBookCatalogue()
	for _, cataLog := range cataLogData {
		if len(cataLog.CateName) > 0 {
			saveData := &models.BookCatalogueCategory{
				BookId:  bookData.Id,
				CatName: cataLog.CateName,
				Sort:    cataLog.Sort,
				AddDate: time.Now().Unix(),
			}
			if catId, err = cataCategoryObj.Save(saveData); err != nil {
				logrus.Infoln("catalogueObj Save error:", err)
				continue
			}
		}
		for _, v := range cataLog.CatalogList {
			logData := &models.BookCatalogue{
				BookId:           bookData.Id,
				CatalogTitle:     v.CatalogTitle,
				CatalogCatgoryId: catId,
				LinkUrl:          v.LinkUrl,
				Sort:             v.Sort,
				AddDate:          time.Now().Unix(),
			}
			if _, err = catalogueObj.Save(logData); err != nil {
				logrus.Infoln("保存目录", v.CatalogTitle, "失败")
			}
		}
	}
	logrus.Infoln("LoadCatalogueData 保存目录信息结束......")
	return
}
