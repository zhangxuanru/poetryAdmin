/*
@Time : 2019/9/5 15:42
@Author : zxr
@File : content
@Software: GoLand
*/
package Action

import (
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/ancient/Parser"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"strings"
	"time"
)

type content struct {
}

func NewContent() *content {
	return &content{}
}

//载入目录和书籍封面信息，获取目录对应的具体内容
func (c *content) LoadBookCatalogue(book *define.BookCatalogue) {
	for _, catalogue := range book.CatalogList {
		go c.procGetContent(catalogue)
		time.Sleep(500 * time.Millisecond)
	}
}

//处理每个目录对应的详情数据
func (c *content) procGetContent(catalogue define.CataLog) {
	var (
		url         string
		bytes       []byte
		err         error
		bookContent *define.BookCatalogueContent
	)
	url = catalogue.LinkUrl
	if strings.Contains(url, "http") == false {
		url = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(url, "/")
	}
	if bytes, err = c.getContentHtml(url); err != nil || len(bytes) == 0 {
		logrus.Infoln("getContentHtml err:", err)
		return
	}
	if bookContent, err = Parser.ParseGuWenContent(bytes); err != nil {
		logrus.Infoln("ParseGuWenContent error:", err)
		return
	}
	if bytes, err = base.GetHtml(bookContent.TranslationUrl); err == nil {
		bookContent.Translation = string(bytes)
	}
	bookContent.ShortCatalogueTitle = catalogue.CatalogTitle
	bookContent.CatalogueLinkUrl = catalogue.LinkUrl

	//保存数据
	sendData := &define.ParseData{
		Data:      bookContent,
		Params:    nil,
		ParseFunc: data.NewBookContentStore().LoadBookContentData,
	}
	data.G_GraspResult.SendParseData(sendData)
	return
}

//获取HTML资源
func (c *content) getContentHtml(url string) (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		dir, _ := os.Getwd()
		file := dir + "/ancient/content.html"
		return tools.ReadFile(file)
	} else {
		bytes, err = base.GetHtml(url)
	}
	return
}
