/*
@Time : 2019/9/2 19:50
@Author : zxr
@File : AncientBookStore
@Software: GoLand
*/
package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/core/define"
)

//保存古籍-书籍封面信息
type AncientBookStore struct {
}

func NewAncientBookStore() *AncientBookStore {
	return &AncientBookStore{}
}

//载入书籍封面信息 明日继续........
func (a *AncientBookStore) LoadBookData(data interface{}, params interface{}) {
	var (
		bookCoverList []*define.GuWenBookCover
		ok            bool
	)
	if bookCoverList, ok = data.([]*define.GuWenBookCover); ok == false {
		logrus.Infoln("LoadBookData data conver to GuWenBookCover error")
		return
	}
	logrus.Infoln("LoadBookData.......")
	logrus.Infoln(bookCoverList)

}
