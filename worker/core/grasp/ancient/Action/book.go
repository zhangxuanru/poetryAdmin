/*
@Time : 2019/9/2 16:18
@Author : zxr
@File : book
@Software: GoLand
*/
package Action

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/core/define"
)

type Book struct {
	categoryHtmlChan chan *define.GuWenCategoryBookHtml
	closeChan        chan bool
}

func NewBook() *Book {
	return &Book{
		categoryHtmlChan: make(chan *define.GuWenCategoryBookHtml, 100),
		closeChan:        make(chan bool),
	}
}

func (b *Book) GrabCategoryBook() {

}

//接收分类下的书籍信息
func (b *Book) ReceiveCategoryBook() {
	var (
		htmlCategory *define.GuWenCategoryBookHtml
		autoClose    bool
	)
	for {
		if len(b.categoryHtmlChan) == 0 && autoClose == true {
			logrus.Infoln("autoClose == true --:")
			goto END
		}
		select {
		case htmlCategory = <-b.categoryHtmlChan:
			logrus.Infoln("htmlCategory:")
			logrus.Infoln(htmlCategory)
			logrus.Infoln("----------------------")
		case <-b.closeChan:
			if len(b.categoryHtmlChan) > 0 {
				logrus.Infoln("categoryHtmlChan还有数据....处理完后才会退出,len:", len(b.categoryHtmlChan))
				autoClose = true
			} else {
				logrus.Infoln("closeChan----")
				goto END
			}
		}
	}
END:
	logrus.Infoln("close ReceiveCategoryBook")
	return
}

//发送抓取的分类书籍信息
func (b *Book) SendCategoryBook(data *define.GuWenCategoryBookHtml) {
	logrus.Infoln("SendCategoryBook...............")
	b.categoryHtmlChan <- data
}

func (b *Book) SendClose(ret bool) {
	b.closeChan <- ret
}
