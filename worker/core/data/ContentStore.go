package data

import (
	"errors"
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"strings"
	"time"
)

//保存诗词正文和赏析注释结果...
type contentStore struct {
	detail *define.PoetryContent
	params define.LinkStr
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
		ok  bool
		err error
	)
	if c.detail, ok = data.(*define.PoetryContent); ok == false {
		logrus.Infoln("LoadPoetryContentData error:data conver define.PoetryContent error")
		return
	}
	if c.params, ok = params.(define.LinkStr); ok == false {
		logrus.Infoln("LoadPoetryContentData error: params conver LinkStr error")
	}
	if err = c.SetAttrData(); err != nil {
		return
	}
	//1.写入诗词表， 先写诗词表，拿到诗词ID 再做下面的处理
	if err = c.SaveContent(); err != nil {
		logrus.Infoln("SaveContent error:", err)
		return
	}
	//2.保存分类
	c.SavePoetryCategory()
	// 3.保存文本内容,翻译表,赏析表
	c.SaveContentNotes()
	//保存到 ES中
	go NewEsStore().SaveContentDetail(c.detail)
	logrus.Infoln("LoadPoetryContentData.............")
}

//设置基本数据
func (c *contentStore) SetAttrData() (err error) {
	var (
		dynastyId int64
		author    models.Author
	)
	if author, err = models.GetAuthorDataByAuthorName(c.detail.Author.AuthorName); err != nil {
		logrus.Infoln("GetAuthorDataByAuthorName error:", err)
		return
	}
	if author.Id > 0 {
		c.detail.Author.AuthorId = int64(author.Id)
		c.detail.Author.DynastyId = author.DynastyId
	} else {
		if dynastyId, err = models.NewDynasty().GetIdBySaveName(c.detail.Author.DynastyName); err == nil {
			c.detail.Author.DynastyId = int(dynastyId)
		}
		authorMod := &models.Author{
			Author:      c.detail.Author.AuthorName,
			DynastyId:   c.detail.Author.DynastyId,
			SourceUrl:   c.detail.Author.AuthorSrcUrl,
			WorksUrl:    c.detail.Author.AuthorContentUrl,
			PhotoUrl:    c.detail.Author.AuthorImgUrl,
			AuthorIntro: c.detail.Author.Introduction,
			PoetryCount: c.detail.Author.AuthorTotalPoetry,
		}
		c.detail.Author.AuthorId, err = models.NewAuthor().SaveAuthor(authorMod)
	}
	if c.detail.Author.AuthorId == 0 {
		logrus.Infoln("content.Author.AuthorId eq 0; err:", err)
	}
	return
}

//写入诗词表和关联表,保存诗词
func (c *contentStore) SaveContent() (err error) {
	var (
		poetryId int64
		notesId  int64
		urlCrc32 uint32
		content  models.Content
	)
	if len(c.detail.SourceUrl) == 0 {
		return errors.New("SaveContent SourceUrl is nil")
	}
	urlCrc32 = tools.Crc32(c.detail.SourceUrl)
	if content, err = models.NewContent().GetContentByCrc32(urlCrc32); err != nil {
		return err
	}
	if len(c.detail.CreativeBackground) > 0 && content.CreatBackId == 0 {
		notesData := models.Notes{
			Title:      "创作背景",
			Content:    c.detail.CreativeBackground,
			Type:       1,
			AddDate:    time.Now().Unix(),
			UpdateDate: time.Now().Unix(),
		}
		notesId, err = models.NewNotes().SaveNotes(&notesData)
	}
	data := &models.Content{
		Title:          strings.TrimSpace(c.detail.Title),
		Content:        strings.TrimSpace(c.detail.Content),
		AuthorId:       c.detail.Author.AuthorId,
		GenreId:        0,
		CreatBackId:    notesId,
		Sort:           c.params.Sort,
		SourceUrl:      c.detail.SourceUrl,
		SourceUrlCrc32: urlCrc32,
		AddDate:        time.Now().Unix(),
		UpdateDate:     time.Now().Unix(),
	}
	if content.Id > 0 {
		data.Id = content.Id
	}
	if poetryId, err = models.NewContent().SaveUpdate(data); err == nil {
		c.detail.PoetryId = poetryId
	}
	return
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
			G_GraspResult.PushError(err, "SaveDetailCategory")
		}
	}
	return
}

//写入文本内容表
func (c *contentStore) SaveContentNotes() {
	for _, cont := range c.detail.Detail {
		if cont.TransId > 0 {
			c.SaveTransNotes(cont) //翻译信息
		}
		if cont.AppRecId > 0 {
			c.SaveRecNotes(cont) //赏析信息
		}
	}
}

//写入翻译表和文本表
func (c *contentStore) SaveTransNotes(data *define.PoetryContentData) {
	var (
		err       error
		transData models.ContentTrans
		notesData *models.Notes
		notesId   int64
	)
	if transData, err = models.NewContentTrans().FindNotesId(int(c.detail.PoetryId), data.TransId); err != nil {
		logrus.Infoln("FindNotesId error:", err)
		return
	}
	data.Content = tools.TrimDivHtml(data.Content)
	notesData = &models.Notes{
		Title:      strings.TrimSpace(data.Title),
		Content:    data.Content,
		PlayUrl:    data.PlayUrl,
		PlaySrcUrl: data.PlaySrcUrl,
		HtmlSrcUrl: data.HtmlSrcUrl,
		Type:       1,
		Introd:     strings.TrimSpace(data.Introd),
		AddDate:    time.Now().Unix(),
		UpdateDate: time.Now().Unix(),
	}
	if transData.Id > 0 {
		notesData.Id = int(transData.NotesId)
	}
	//写|更新详情表
	if notesId, err = models.NewNotes().SaveNotes(notesData); err != nil {
		logrus.Infoln("SaveNotes error:", err)
		return
	}
	//上传mp3
	if len(data.PlayUrl) > 0 {
		go NewAuthorStore().UploadMp3(notesId, data.PlayUrl)
	}
	if transData.Id > 0 {
		return
	}
	//写翻译表
	transData = models.ContentTrans{
		PoetryId:   int(c.detail.PoetryId),
		TransId:    data.TransId,
		NotesId:    notesId,
		Sort:       data.Sort,
		AddDate:    time.Now().Unix(),
		UpdateDate: time.Now().Unix(),
	}
	if _, err = models.NewContentTrans().InsertTrans(&transData); err != nil {
		logrus.Infoln("InsertTrans error:", err)
		G_GraspResult.PushError(err, "InsertTrans error")
		return
	}
}

//写入赏析信息表和文本表
func (c *contentStore) SaveRecNotes(data *define.PoetryContentData) {
	var (
		err       error
		recData   models.ContentRec
		notesData *models.Notes
		notesId   int64
	)
	if recData, err = models.NewContentRec().FindNotesId(int(c.detail.PoetryId), data.AppRecId); err != nil {
		logrus.Infoln("ContentRec FindNotesId error:", err)
		return
	}
	data.Content = tools.TrimDivHtml(data.Content)
	notesData = &models.Notes{
		Title:      strings.TrimSpace(data.Title),
		Content:    data.Content,
		PlayUrl:    data.PlayUrl,
		PlaySrcUrl: data.PlaySrcUrl,
		HtmlSrcUrl: data.HtmlSrcUrl,
		Type:       1,
		Introd:     strings.TrimSpace(data.Introd),
		AddDate:    time.Now().Unix(),
		UpdateDate: time.Now().Unix(),
	}
	if recData.Id > 0 {
		notesData.Id = int(recData.NotesId)
	}
	//写|更新详情表
	if notesId, err = models.NewNotes().SaveNotes(notesData); err != nil {
		logrus.Infoln("SaveNotes error:", err)
		return
	}
	//上传mp3
	if len(data.PlayUrl) > 0 {
		go NewAuthorStore().UploadMp3(notesId, data.PlayUrl)
	}
	if recData.Id > 0 {
		return
	}
	//写翻译表
	recData = models.ContentRec{
		PoetryId:   int(c.detail.PoetryId),
		ApprecId:   data.AppRecId,
		NotesId:    notesId,
		Sort:       data.Sort,
		AddDate:    time.Now().Unix(),
		UpdateDate: time.Now().Unix(),
	}
	if _, err = models.NewContentRec().InsertRec(&recData); err != nil {
		logrus.Infoln("InsertTrans error:", err)
		G_GraspResult.PushError(err, "InsertTrans error")
		return
	}
}
