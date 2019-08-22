package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/core/define"
)

//保存作者信息
type AuthorStore struct {
}

func NewAuthorStore() *AuthorStore {
	return new(AuthorStore)
}

//载入作者信息并保存数据
func (c *AuthorStore) LoadAuthorData(data interface{}, params interface{}) {
	detail := data.(*define.PoetryAuthorDetail)
	if len(detail.AuthorName) == 0 {
		logrus.Infoln("LoadAuthorData AuthorName is nil")
		return
	}
	for k, v := range detail.Data {
		logrus.Infoln(k, ":", v)
	}
	logrus.Infof("%+v", detail)
}
