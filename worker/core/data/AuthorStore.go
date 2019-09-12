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
	return &AuthorStore{}
}

//载入作者信息和资料信息并保存数据
func (a *AuthorStore) LoadAuthorData(data interface{}, params interface{}) {
	var err error
	detail := data.(*define.PoetryAuthorDetail)
	if len(detail.AuthorName) == 0 {
		logrus.Infoln("LoadAuthorData AuthorName is nil")
		return
	}
	//1.获取朝代信息poetry_dynasty
	detail.DynastyName = strings.TrimSpace(detail.DynastyName)
	if len(detail.DynastyName) > 0 {
		if dynastyId, err := models.NewDynasty().GetIdBySaveName(detail.DynastyName); err == nil {
			detail.DynastyId = int(dynastyId)
		}
	}
	//2.更新作者表poetry_author
	if detail.AuthorId, err = a.UpdateAuthor(detail); err != nil || detail.AuthorId == 0 {
		logrus.Infoln("更新作者信息错误 err:", err)
		return
	}
	//上传头像
	go a.UploadFile(detail)

	//3.写入poetry_detail_notes表
	for k, val := range detail.Data {
		var (
			notesId    int64
			authorData models.AuthorData
		)
		//查询author_data判断是否已经存在
		authorData, _ = models.NewAuthorData().GetAuthorDataByDataId(detail.AuthorId, int64(val.DataId))
		if authorData.NotesId > 0 {
			val.Id = authorData.NotesId
		}
		val.Sort = k
		if notesId, err = NewNotesStore().SaveNotes(val); err != nil {
			continue
		}
		//上传mp3
		go a.UploadMp3(notesId, val.PlayUrl)

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

//执行数据库操作 更新作者信息
func (a *AuthorStore) UpdateAuthor(detail *define.PoetryAuthorDetail) (id int64, err error) {
	author := &models.Author{
		Author:      detail.AuthorName,
		SourceUrl:   detail.AuthorSrcUrl,
		WorksUrl:    detail.AuthorContentUrl,
		DynastyId:   detail.DynastyId,
		PhotoUrl:    detail.AuthorImgUrl,
		AuthorIntro: detail.Introduction,
		PoetryCount: detail.AuthorTotalPoetry,
	}
	if len(detail.AuthorImgFileName) > 0 {
		author.PhotoFileName = detail.AuthorImgFileName
	}
	id, err = models.NewAuthor().UpdateAuthor(author)
	return
}

//上传头像
func (a *AuthorStore) UploadFile(detail *define.PoetryAuthorDetail) {
	var (
		fileName string
		err      error
	)
	if len(detail.AuthorImgUrl) > 0 {
		if fileName, err = NewUploadStore().Upload(detail.AuthorImgUrl); err != nil {
			logrus.Infoln("upload image error:", err)
		} else {
			author := &models.Author{
				Author:        detail.AuthorName,
				PhotoFileName: fileName,
			}
			_, _ = models.NewAuthor().UpdateAuthorPhoto(author)
		}
	}
	return
}

//上传mp3
func (a *AuthorStore) UploadMp3(id int64, mp3Url string) {
	var (
		fileName string
		err      error
	)
	if len(mp3Url) > 0 {
		if fileName, err = NewUploadStore().Upload(mp3Url); err != nil {
			logrus.Infoln(mp3Url, "upload error:", err)
			return
		}
		//更新MP3地址
		data := &define.ContentData{
			Id:       int(id),
			FileName: fileName,
		}
		if _, err = NewNotesStore().UpdateMp3Path(data); err != nil {
			logrus.Infoln("UpdateMp3Path error:", err)
		}
	}
	return
}
