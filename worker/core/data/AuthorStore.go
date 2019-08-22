package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/core/define"
	"strings"
)

//保存作者信息
type AuthorStore struct {
}

func NewAuthorStore() *AuthorStore {
	return new(AuthorStore)
}

//载入作者信息并保存数据
func (a *AuthorStore) LoadAuthorData(data interface{}, params interface{}) {
	detail := data.(*define.PoetryAuthorDetail)
	if len(detail.AuthorName) == 0 {
		logrus.Infoln("LoadAuthorData AuthorName is nil")
		return
	}
	//1.获取朝代信息poetry_dynasty
	detail.DynastyName = strings.TrimSpace(detail.DynastyName)
	if dynastyId, err := models.NewDynasty().GetIdBySaveName(detail.DynastyName); err == nil {
		detail.DynastyId = int(dynastyId)
	}
	//上传作者头像
	if len(detail.AuthorImgUrl) > 0 {
		if fileName, err := NewImgStore().UploadImg(detail.AuthorImgUrl); err != nil {
			logrus.Infoln("upload image error:", err)
		} else {
			detail.AuthorImgFileName = fileName
		}
	}
	//2.更新作者表poetry_author

	//3.写入poetry_detail_notes表
	//4.写入poetry_author_data表
	for k, v := range detail.Data {
		logrus.Infoln(k, ":", v)
	}
	logrus.Infof("%+v", detail)
}
