/*
@Time : 2019/9/2 19:50
@Author : zxr
@File : AncientBookStore
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

//保存古籍-书籍封面信息
type AncientBookStore struct {
}

func NewAncientBookStore() *AncientBookStore {
	return &AncientBookStore{}
}

//载入书籍封面信息
func (a *AncientBookStore) LoadBookData(data interface{}, params interface{}) {
	var (
		bookCoverList []*define.GuWenBookCover
		category      *define.GuWenCategory
		categoryData  models.AncientCategory
		ok            bool
		bookId        int64
		err           error
		bookObj       *models.AncientBook
		categoryObj   *models.AncientCategory
	)
	if bookCoverList, ok = data.([]*define.GuWenBookCover); ok == false {
		logrus.Infoln("LoadBookData data conver to GuWenBookCover error")
		return
	}
	if category, ok = params.(*define.GuWenCategory); ok == false {
		logrus.Infoln("LoadBookData params conver to GuWenCategory error")
		return
	}
	bookObj = models.NewAncientBook()
	categoryObj = models.NewAncientCategory()
	for _, book := range bookCoverList {
		if categoryData, err = categoryObj.GetCategoryDataByName(category.CategoryName); err != nil {
			logrus.Infoln("GetCategoryDataByName err:", err)
			continue
		}
		if categoryData.Id == 0 {
			logrus.Infoln("分类不存在重新插入.... ")
			categoryMd := &models.AncientCategory{
				CatName: category.CategoryName,
				SrcUrl:  category.LinkUrl,
				Sort:    category.Sort,
			}
			if cId, err := categoryObj.SaveCategory(categoryMd); err != nil {
				logrus.Infoln("SaveCategory error:", err)
				continue
			} else {
				categoryData.Id = int(cId)
			}
		}
		bookMod := &models.AncientBook{
			CatId:            categoryData.Id,
			BookTitle:        book.Title,
			BookIntroduction: book.Introduction,
			LinkUrl:          book.LinkUrl,
			LinkUrlCrc32:     tools.Crc32(book.LinkUrl),
			SongUrl:          book.SongUrl,
			SongSrcUrl:       book.SongSrcUrl,
			FamousTotal:      book.FamousTotal,
			CoverChart:       book.CoverChart,
			AddDate:          time.Now().Unix(),
		}
		if bookId, err = bookObj.SaveBook(bookMod); err != nil {
			logrus.Infoln("SaveBook error:", err)
			continue
		}
		if bookId == 0 {
			logrus.Infoln("bookId is nil")
			continue
		}
		go a.UploadMp3(bookId, book.SongUrl)
		go a.UploadCover(bookId, book.CoverChart)
	}
	logrus.Infoln("LoadBookData end .......")
	return
}

//更新声音文件
func (a *AncientBookStore) UploadMp3(id int64, mp3Url string) {
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
	if _, err = models.NewAncientBook().UpdateSongPath(id, fileName); err != nil {
		logrus.Infoln("UpdateSongPath error:", err)
		return
	}
	return
}

//更新封面图片地址
func (a *AncientBookStore) UploadCover(id int64, imgUrl string) {
	var (
		fileName string
		err      error
	)
	if id == 0 || len(imgUrl) == 0 {
		return
	}
	if fileName, err = NewUploadStore().Upload(imgUrl); err != nil {
		logrus.Infoln("upload image error:", err)
		return
	}
	if _, err = models.NewAncientBook().UpdateCoverPath(id, fileName); err != nil {
		logrus.Infoln("UpdateCoverPath  error:", err)
		return
	}
	return
}
