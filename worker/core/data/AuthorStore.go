package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/core/define"
	"strings"
	"time"
)

//保存作者信息
type AuthorStore struct {
}

func NewAuthorStore() *AuthorStore {
	return new(AuthorStore)
}

//载入作者信息并保存数据 这里保存数据慢， 下周需要优化.......
func (a *AuthorStore) LoadAuthorData(data interface{}, params interface{}) {
	var (
		err error
	)
	detail := data.(*define.PoetryAuthorDetail)
	if len(detail.AuthorName) == 0 {
		logrus.Infoln("LoadAuthorData AuthorName is nil")
		return
	}
	//1.获取朝代信息poetry_dynasty
	a.getAttr(detail)
	//2.更新作者表poetry_author
	if detail.AuthorId, err = a.UpdateAuthor(detail); err != nil || detail.AuthorId == 0 {
		logrus.Infoln("更新作者信息错误 err:", err)
		return
	}
	//3.写入poetry_detail_notes表
	for _, val := range detail.Data {
		var (
			notesId    int64
			authorData models.AuthorData
		)
		//上传mp3
		if len(val.PlayUrl) > 0 {
			if val.FileName, err = NewUploadStore().Upload(val.PlayUrl); err != nil {
				logrus.Infoln(val.PlayUrl, "upload error:", err)
				val.FileName = ""
			}
		}
		//查询author_data判断是否已经存在
		authorData, _ = models.NewAuthorData().GetNotesByDataId(detail.AuthorId, int64(val.DataId))
		if authorData.NotesId > 0 {
			val.Id = authorData.NotesId
		}
		if notesId, err = NewNotesStore().SaveNotes(val); err != nil {
			continue
		}
		//4.写入poetry_author_data表
		info := &models.AuthorData{
			AuthorId:   detail.AuthorId,
			DataId:     val.DataId,
			NotesId:    int(notesId),
			AddDate:    time.Now().Unix(),
			UpdateDate: time.Now().Unix(),
		}
		if authorData.Id > 0 {
			info.Id = authorData.Id
		}
		if _, err = models.NewAuthorData().SaveAuthorData(info); err != nil {
			logrus.Infoln("SaveAuthorData error:", err)
		}
	}
	//保存作者数据到ES中
	NewEsStore().SaveAuthorData(detail)
	logrus.Infoln("LoadAuthorData end .....")
}

//获取朝代ID和图片上传路径
func (a *AuthorStore) getAttr(detail *define.PoetryAuthorDetail) {
	detail.DynastyName = strings.TrimSpace(detail.DynastyName)
	if dynastyId, err := models.NewDynasty().GetIdBySaveName(detail.DynastyName); err == nil {
		detail.DynastyId = int(dynastyId)
	}
	//上传作者头像
	if len(detail.AuthorImgUrl) > 0 {
		if fileName, err := NewUploadStore().Upload(detail.AuthorImgUrl); err != nil {
			logrus.Infoln("upload image error:", err)
		} else {
			detail.AuthorImgFileName = fileName
		}
	}
}

//更新作者信息
func (a *AuthorStore) UpdateAuthor(detail *define.PoetryAuthorDetail) (id int64, err error) {
	author := &models.Author{
		Author:        detail.AuthorName,
		SourceUrl:     detail.AuthorSrcUrl,
		WorksUrl:      detail.AuthorContentUrl,
		DynastyId:     detail.DynastyId,
		PhotoUrl:      detail.AuthorImgUrl,
		PhotoFileName: detail.AuthorImgFileName,
		AuthorIntro:   detail.Introduction,
		PoetryCount:   detail.AuthorTotalPoetry,
	}
	id, err = models.NewAuthor().UpdateAuthor(author)
	return
}
