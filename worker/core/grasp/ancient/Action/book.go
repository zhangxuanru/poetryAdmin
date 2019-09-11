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

//书籍封面信息
type Book struct {
	categoryHtmlChan chan *define.GuWenCategoryBookHtml
	closeChan        chan bool
}

func NewBook() *Book {
	return &Book{
		categoryHtmlChan: make(chan *define.GuWenCategoryBookHtml),
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
		if len(b.categoryHtmlChan) == 0 && htmlCategory == nil && autoClose == true {
			goto END
		}
		select {
		case htmlCategory = <-b.categoryHtmlChan:
			if bookCoverList, err = Parser.ParseGuWenCategoryBook(htmlCategory.Html); err != nil {
				logrus.Infoln("ParseGuWenCategoryBook error:", err)
				continue
			}
			//logrus.Infoln("---------------------")
			//logrus.Infoln("url:", htmlCategory.LinkUrl)
			//logrus.Infof("GuWenCategory:%+v", htmlCategory.GuWenCategory)
			//logrus.Infof("bookCoverList:%+v\n", bookCoverList)
			for _, r := range bookCoverList {
				logrus.Infof("%v\n", r)
				if r.LinkUrl == "" {
					logrus.Infoln("url=nil", r, "--topUrl:", htmlCategory.LinkUrl)
				}
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
	/*
		sendData := &define.ParseData{
			Data:      bookCoverList,
			Params:    category,
			ParseFunc: data.NewAncientBookStore().LoadBookData,
		}
		data.G_GraspResult.SendParseData(sendData)
	*/
	data.NewAncientBookStore().LoadBookData(bookCoverList, category)
	//发送书籍目录页请求
	NewCataLog().LoadBookCoverList(bookCoverList, category)
}

//发送抓取的分类书籍信息
func (b *Book) SendCategoryBook(data *define.GuWenCategoryBookHtml) {
	logrus.Infof("SendCategoryBook data:%+v\n", data.GuWenCategory)
	b.categoryHtmlChan <- data
}

func (b *Book) SendClose(ret bool) {
	b.closeChan <- ret
}
