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

//载入名句主题数据
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

//保存名句-主题下的分类数据
func (f *FamousStorage) LoadFamousCategoryData(data interface{}, params interface{}) {
	var (
		themeCategory *define.ThemeCategory
		categorys     models.Category
		catId         int64
		err           error
		ok            bool
	)
	if themeCategory, ok = data.(*define.ThemeCategory); ok == false {
		logrus.Infoln("LoadFamousCategoryData data conver to ThemeCategory error")
		return
	}
	if categorys, err = models.GetDataByCateName(themeCategory.Title, 2); err != nil {
		logrus.Infoln("GetDataByCateName err:", err)
		G_GraspResult.PushError(err)
		return
	}
	catId = int64(categorys.Id)
	if catId == 0 {
		categorys = models.Category{
			CatName:        themeCategory.Title,
			SourceUrl:      themeCategory.LinkUrl,
			SourceUrlCrc32: tools.Crc32(themeCategory.LinkUrl),
			ShowPosition:   2,
			AddDate:        time.Now().Unix(),
		}
		if catId, err = models.SaveCategoryData(&categorys); err != nil {
			logrus.Infoln("SaveCategoryData err:", err)
			return
		}
	}
	if catId == 0 {
		logrus.Infoln("LoadFamousCategoryData catId is zero")
		return
	}
	for _, classify := range themeCategory.ClassifyList {
		category := &models.Category{
			CatName:        classify.Title,
			SourceUrl:      classify.LinkUrl,
			SourceUrlCrc32: tools.Crc32(classify.LinkUrl),
			ShowPosition:   2,
			Pid:            int(catId),
			AddDate:        time.Now().Unix(),
		}
		if catData, err := models.GetDataByCateName(category.CatName, category.ShowPosition); err != nil || catData.Id > 0 {
			logrus.Infoln("GetDataByCateName err:", err, "ID:", catData.Id)
			continue
		}
		if _, err = models.SaveCategoryData(category); err != nil {
			logrus.Infoln("SaveCategoryData err:", err)
			continue
		}
	}
}

//保存名句-分类下的名句列表
func (f *FamousStorage) LoadClassifyContentData(data interface{}, params interface{}) {
	var (
		ok           bool
		classify     *define.Classify
		category     models.Category
		content      define.Content
		sentenceData *models.FamousSentence
		sentenceObj  *models.FamousSentence
		err          error
		id           int64
	)
	if classify, ok = data.(*define.Classify); ok == false {
		logrus.Infoln("LoadClassifyContentData data conver to Classify error")
		return
	}
	if category, err = models.GetDataByCateNameAndPid(classify.ThemeTitle, 2, 0); err != nil || category.Id == 0 {
		logrus.Infoln("GetDataByCateNameAndPid err:", err)
		return
	}
	if category, err = models.GetDataByCateNameAndPid(classify.Title, 2, category.Id); err != nil || category.Id == 0 {
		logrus.Infoln("GetDataByCateNameAndPid get ", classify.Title, "-err:", err)
		return
	}
	//写入 poetry_famous_sentence
	sentenceObj = models.NewFamousSentence()
	for _, content = range classify.ContentList {
		title := content.PoetryTitle
		if len(title) == 0 {
			title = content.PoetryText
		}
		sentenceData = &models.FamousSentence{
			CatId:        category.Id,
			Content:      content.Text,
			ContentCrc32: tools.Crc32(content.Text),
			PoetryTitle:  title,
			Sort:         content.Sort,
			SourceUrl:    content.PoetryLink,
			SourceCrc32:  tools.Crc32(content.PoetryLink),
			AddDate:      time.Now().Unix(),
		}
		if sentence, err := sentenceObj.GetDataByCrc32(sentenceData.SourceCrc32, sentenceData.ContentCrc32); err != nil || sentence.Id > 0 {
			logrus.Infoln("sentenceObj GetDataByCrc32 err:", err, "-Id:", sentence.Id)
			continue
		}
		if id, err = sentenceObj.Save(sentenceData); err != nil {
			logrus.Infoln("sentenceObj save error:", err)
			continue
		}
		if id == 0 {
			logrus.Infoln("保存--", title, "--名句失败.....")
			continue
		}
	}
}
