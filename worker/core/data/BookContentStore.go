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

//保存古籍-详情内容
type BookContentStore struct {
}

func NewBookContentStore() *BookContentStore {
	return &BookContentStore{}
}

//载入详情内容数据
func (b *BookContentStore) LoadBookContentData(data interface{}, params interface{}) {
	var (
		bookContent   *define.BookCatalogueContent
		saveData      *models.BookContent
		catalogueData models.BookCatalogue
		content       models.BookContent
		authorData    *models.AncientAuthor
		authorId      int64
		id            int64
		err           error
		ok            bool
	)
	if bookContent, ok = data.(*define.BookCatalogueContent); ok == false {
		logrus.Infoln("LoadBookContentData data conver to BookContent error")
		return
	}
	//查询目录ID和书籍ID
	if catalogueData, err = models.NewBookCatalogue().GetDataByTitleAndUrl(bookContent.ShortCatalogueTitle, bookContent.CatalogueLinkUrl); err != nil {
		logrus.Infoln("GetDataByTitleAndUrl error:", err)
		return
	}
	if catalogueData.Id == 0 {
		logrus.Infoln("目录内容不存在")
		return
	}
	if content, err = models.NewBookContent().GetContentByLogId(catalogueData.Id); err != nil {
		logrus.Infoln("GetContentByLogId error:", err)
		return
	}
	if content.Id > 0 {
		logrus.Infoln(bookContent.CatalogueTitle, "已存在")
		return
	}
	//查询作者是否存在
	authorData = &models.AncientAuthor{
		AuthorName: bookContent.AuthorName,
		SourceUrl:  bookContent.AuthorLinkUrl,
	}
	if authorId, err = models.NewAncientAuthor().SaveAuthor(authorData); err != nil {
		logrus.Infoln("保存作者信息失败......")
	}
	saveData = &models.BookContent{
		BookId:         catalogueData.BookId,
		CatalogId:      catalogueData.Id,
		Content:        bookContent.Content,
		Translation:    bookContent.Translation,
		TranslationId:  bookContent.TranslationId,
		TranslationUrl: bookContent.TranslationUrl,
		AuthorId:       authorId,
		SongUrl:        bookContent.SongUrl,
		AddDate:        time.Now().Unix(),
	}
	if id, err = models.NewBookContent().Save(saveData); err != nil {
		logrus.Infoln("BookContent Save error:", err)
		return
	}
	go b.UploadMp3(id, bookContent.SongUrl)
	logrus.Infoln("LoadBookContentData 保存", bookContent.CatalogueTitle, "详情信息结束......")
	return
}

//更新声音文件
func (b *BookContentStore) UploadMp3(id int64, mp3Url string) {
	var (
		fileName string
		err      error
	)
	if len(mp3Url) == 0 || id == 0 {
		return
	}
	if fileName, err = NewUploadStore().Upload(mp3Url); err != nil {
		logrus.Infoln(mp3Url, "upload error:", err)
		return
	}
	if _, err = models.NewBookContent().UpdateSongPath(id, fileName); err != nil {
		logrus.Infoln("UpdateSongPath error:", err)
		return
	}
	return
}
