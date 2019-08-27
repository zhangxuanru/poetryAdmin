package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"time"
)

//保存诗词正文和赏析注释结果...
type contentStore struct {
	detail *define.PoetryContent
}

func NewContentStore() *contentStore {
	return new(contentStore)
}

//保存诗词详情信息与资料
func (c *contentStore) LoadPoetryContentData(data interface{}, params interface{}) {
	/*
	 1.写入poetry_content诗词内容
	 2.写入poetry_detail_category 诗词分类表，
	 3.写入poetry_detail_notes 文本内容表
	 4.写入poetry_content_trans 翻译表
	 5.写入poetry_content_apprec赏析表
	 6.更新poetry_content_relation表
	*/
	var (
		content *define.PoetryContent
		ok      bool
	)
	if content, ok = data.(*define.PoetryContent); ok == false {
		logrus.Infoln("LoadPoetryContentData error:data conver define.PoetryContent error")
		return
	}
	c.SetAttrData(content)
	c.detail = content
	//1.写入诗词表， 先写诗词表，拿到  诗词ID 再做下面的处理

	//2.保存分类
	c.SavePoetryCategory()

	//保存到 ES中
	NewEsStore().SaveContentDetail(content)

	logrus.Infoln("LoadPoetryContentData.............")
	logrus.Infof("%+v\n", data)
	logrus.Infof("%+v\n", params)
}

//设置基本数据
func (c *contentStore) SetAttrData(content *define.PoetryContent) {
	var (
		dynastyId  int64
		err        error
		author     models.Author
		contentRow models.Content
	)
	if author, err = models.GetAuthorDataByAuthorName(content.Author.AuthorName); err != nil {
		logrus.Infoln("GetAuthorDataByAuthorName error:", err)
		return
	}
	if author.Id > 0 {
		content.Author.AuthorId = int64(author.Id)
		content.Author.DynastyId = author.DynastyId
	} else {
		if dynastyId, err = models.NewDynasty().GetIdBySaveName(content.Author.DynastyName); err == nil {
			content.Author.DynastyId = int(dynastyId)
		}
		authorMod := &models.Author{
			Author:      content.Author.AuthorName,
			DynastyId:   content.Author.DynastyId,
			SourceUrl:   content.Author.AuthorSrcUrl,
			WorksUrl:    content.Author.AuthorContentUrl,
			PhotoUrl:    content.Author.AuthorImgUrl,
			AuthorIntro: content.Author.Introduction,
			PoetryCount: content.Author.AuthorTotalPoetry,
		}
		content.Author.AuthorId, err = models.NewAuthor().SaveAuthor(authorMod)
	}
	if content.Author.AuthorId == 0 {
		logrus.Infoln("content.Author.AuthorId eq 0; err:", err)
		return
	}
	if contentRow, err = models.NewContent().GetByTitleAuthorId(content.Title, content.Author.AuthorId); err != nil {
		logrus.Infoln("GetByTitleAuthorId err:", err)
		return
	}
	content.PoetryId = int64(contentRow.Id)
}

//保存诗词分类
func (c *contentStore) SavePoetryCategory() {
	var (
		categoryList []*define.TextHrefFormat
		categorys    models.Category
		categoryId   int64
		err          error
	)
	categoryList = c.detail.CategoryList
	if len(categoryList) == 0 {
		logrus.Infoln("err: categoryList is nil")
		return
	}
	for _, cate := range categoryList {
		urlCrc := tools.Crc32(cate.Href)
		if categorys, err = models.GetDataByCrcAndCateName(urlCrc, cate.Text, int(cate.ShowPosition)); err != nil {
			logrus.Infoln("GetDataByCrcAndCateName err:", err)
			continue
		}
		categoryId = int64(categorys.Id)
		if categorys.Id == 0 {
			category := &models.Category{
				CatName:        cate.Text,
				SourceUrl:      cate.Href,
				SourceUrlCrc32: tools.Crc32(cate.Href),
				ShowPosition:   int(cate.ShowPosition),
			}
			categoryId, err = models.SaveCategoryData(category)
		}
		if categoryId == 0 {
			logrus.Infoln("err : categoryId is nil err:", err)
			continue
		}
		//save detail_category
		category := &models.DetailCategory{
			PoetryId:   int(c.detail.PoetryId),
			CategoryId: int(categoryId),
			UpdateTime: time.Now().Unix(),
		}
		if _, err = models.NewDetailCategory().SaveDetailCategory(category); err != nil {
			logrus.Infoln("SaveDetailCategory err:", err)
		}
	}
	return
}
