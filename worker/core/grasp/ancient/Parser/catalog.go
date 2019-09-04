/*
@Time : 2019/9/4 17:36
@Author : zxr
@File : catalog
@Software: GoLand
*/
package Parser

import (
	"github.com/PuerkitoBio/goquery"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
)

//解析分类下书籍目录信息
//https://so.gushiwen.org/guwen/book_24.aspx
func ParseGuWenCataLog(html []byte) (catalogData []*define.CataLogData, err error) {
	var (
		query *goquery.Document
	)
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	query.Find(".sons>.bookcont").Each(func(i int, selection *goquery.Selection) {
		var data define.CataLogData
		data.CateName = selection.Find(".bookMl>strong").Text()
		data.Sort = i
		selection.Find("span>a").Each(func(i int, selection *goquery.Selection) {
			var catalog define.CataLog
			catalog.LinkUrl, _ = selection.Attr("href")
			catalog.CatalogTitle = selection.Text()
			catalog.Sort = i
			data.CatalogList = append(data.CatalogList, catalog)
		})
		catalogData = append(catalogData, &data)
	})
	return
}
