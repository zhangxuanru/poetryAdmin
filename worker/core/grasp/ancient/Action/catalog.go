/*
@Time : 2019/9/4 17:12
@Author : zxr
@File : catalog
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
	"time"
)

//书籍目录相关
type CataLog struct {
}

func NewCataLog() *CataLog {
	return &CataLog{}
}

//载入书籍封面信息，执行抓取目录数据
func (c *CataLog) LoadBookCoverList(bookCoverList []*define.GuWenBookCover, category *define.GuWenCategory) {
	for _, book := range bookCoverList {
		book.GuWenCategory.CategoryName = category.CategoryName
		book.GuWenCategory.LinkUrl = category.LinkUrl
		go c.procBookSource(book)
		time.Sleep(1 * time.Second)
	}
}

//获取html并解析保存，发送详情页请求
func (c *CataLog) procBookSource(book *define.GuWenBookCover) {
	var (
		html          []byte
		err           error
		catalogData   []*define.CataLogData
		bookCatalogue define.BookCatalogue
	)
	if len(book.LinkUrl) == 0 {
		logrus.Infof("book.LinkUrl is nil, book:%+v\n", book)
	}
	if html, err = c.GetCataLogHtml(book.LinkUrl); err != nil {
		logrus.Infoln("getCataLogHtml error:", err)
		return
	}
	if catalogData, err = Parser.ParseGuWenCataLog(html); err != nil {
		logrus.Infoln("ParseGuWenCataLog error:", err)
		return
	}
	//保存数据
	/*
		sendData := &define.ParseData{
			Data:      catalogData,
			Params:    book,
			ParseFunc: data.NewBookCatalogueStore().LoadCatalogueData,
		}
		data.G_GraspResult.SendParseData(sendData)
	*/

	data.NewBookCatalogueStore().LoadCatalogueData(catalogData, book)

	//发送书籍详情内容的请求
	bookCatalogue.BookTitle = book.Title
	bookCatalogue.BookLinkUrl = book.LinkUrl
	for _, cata := range catalogData {
		bookCatalogue.CatalogList = append(bookCatalogue.CatalogList, cata.CatalogList...)
	}
	NewContent().LoadBookCatalogue(&bookCatalogue)
}

//加载页面内容
func (c *CataLog) GetCataLogHtml(url string) (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		dir, _ := os.Getwd()
		file := dir + "/ancient/catalog.html"
		return tools.ReadFile(file)
	} else {
		bytes, err = base.GetHtml(url)
	}
	return
}
