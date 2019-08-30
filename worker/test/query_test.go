/*
@Time : 2019/8/29 20:31
@Author : zxr
@File : testquery
@Software: GoLand
*/
package test

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/grasp/poetry/Helper"
	"testing"
)

func TestCategoryHtml(t *testing.T) {
	bytes, _ := getHtml()
	query, _ := tools.NewDocumentFromReader(string(bytes))
	author := Helper.GetAuthorData(query)

	logrus.Infof("%+v", author)
}

func getHtml() (bytes []byte, err error) {
	dir, _ := os.Getwd()
	file := dir + "/temp/content1.html"
	if ret, _ := tools.PathExists(file); ret == true {
		return tools.ReadFile(file)
	}
	return nil, errors.New(file + "file is not exists")
}
