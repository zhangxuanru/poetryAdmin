/*
@Time : 2019/9/2 16:18
@Author : zxr
@File : book
@Software: GoLand
*/
package Action

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/ancient/Parser"
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

//接收分类下的书籍信息
func (b *Book) ReceiveCategoryBook() {
	var (
		htmlCategory  *define.GuWenCategoryBookHtml
		bookCoverList []*define.GuWenBookCover
		err           error
		autoClose     bool
	)
	for {
		if len(b.categoryHtmlChan) == 0 && autoClose == true {
			goto END
		}
		select {
		case htmlCategory = <-b.categoryHtmlChan:
			if bookCoverList, err = Parser.ParseGuWenCategoryBook(htmlCategory.Html); err != nil {
				logrus.Infoln("ParseGuWenCategoryBook error:", err)
				continue
			}
			b.SaveAndSendBookCover(bookCoverList, &htmlCategory.GuWenCategory)
		case <-b.closeChan:
			if len(b.categoryHtmlChan) > 0 {
				logrus.Infoln("categoryHtmlChan还有数据....处理完后才会退出")
				autoClose = true
			} else {
				goto END
			}
		}
	}
END:
	logrus.Infoln("close ReceiveCategoryBook")
	return
}

//保存书籍封面信息并且发送书籍详情的请求
func (b *Book) SaveAndSendBookCover(bookCoverList []*define.GuWenBookCover, category *define.GuWenCategory) {
	sendData := &define.ParseData{
		Data:      bookCoverList,
		Params:    category,
		ParseFunc: data.NewAncientBookStore().LoadBookData,
	}
	data.G_GraspResult.SendParseData(sendData)
	//发送书籍详情页请求

}

//发送抓取的分类书籍信息
func (b *Book) SendCategoryBook(data *define.GuWenCategoryBookHtml) {
	b.categoryHtmlChan <- data
}

func (b *Book) SendClose(ret bool) {
	b.closeChan <- ret
}
