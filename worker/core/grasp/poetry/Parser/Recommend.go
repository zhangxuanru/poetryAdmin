/*
@Time : 2019/9/16 16:22
@Author : zxr
@File : Recommend
@Software: GoLand
*/
package Parser

import (
	"github.com/PuerkitoBio/goquery"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
)

//解析推荐数据
func ParseRecommendHtml(html []byte) (contentList []define.PoetryContent, page define.ContentPage, err error) {
	var query *goquery.Document
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	query.Find(".left>.sons").Each(func(i int, selection *goquery.Selection) {
		var content define.PoetryContent
		content.Title = selection.Find(".cont").Eq(0).Find("p>a>b").Text()
		content.SourceUrl, _ = selection.Find(".cont").Eq(0).Find("p>a").Attr("href")
		content.Content, _ = selection.Find(".contson").Html()
		content.Sort = i
		contentList = append(contentList, content)
	})
	page.TotalPage = 10
	page.NextPageUrl, _ = query.Find(".left>#FromPage>.pagesright>.amore").Attr("href")
	if len(page.NextPageUrl) > 0 {
		page.IsNextPage = true
	}
	return
}
