/*
@Time : 2019/9/9 18:21
@Author : zxr
@File : content
@Software: GoLand
*/
package Parser

import (
	"github.com/PuerkitoBio/goquery"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"strconv"
	"strings"
)

//解析分类下的详情数据
func ParseFamousContent(html []byte) (contentList []define.Content, page define.ContentPage, err error) {
	var (
		query   *goquery.Document
		content define.Content
	)
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	query.Find(".left>.sons>.cont").Each(func(i int, selection *goquery.Selection) {
		content.Text = selection.Find("a").Eq(0).Text()
		content.LinkUrl, _ = selection.Find("a").Eq(0).Attr("href")
		content.PoetryText = selection.Find("a").Eq(1).Text()
		content.PoetryLink, _ = selection.Find("a").Eq(1).Attr("href")
		content.Sort = i
		if len(content.PoetryLink) > 0 && strings.Contains(content.PoetryLink, "http") == false {
			content.PoetryLink = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(content.PoetryLink, "/")
		}
		if len(content.PoetryText) > 0 {
			content.PoetryText = strings.TrimSpace(content.PoetryText)
			if index := strings.Index(content.PoetryText, "《"); index > 0 {
				content.AuthorName = content.PoetryText[:index]
				content.PoetryTitle = content.PoetryText[index:]
				content.PoetryTitle = strings.TrimLeft(content.PoetryTitle, "《")
				content.PoetryTitle = strings.TrimRight(content.PoetryTitle, "》")
			}
		}
		contentList = append(contentList, content)
	})
	page.IsNextPage = false
	page.NextPageUrl, _ = query.Find("#FromPage>.pagesright>a").Eq(0).Attr("href")
	if len(page.NextPageUrl) > 0 {
		page.NextPageUrl = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(page.NextPageUrl, "/")
		page.IsNextPage = true
	}
	totalPageStr := query.Find("#FromPage>.pagesright>span").Text()
	if len(totalPageStr) > 0 {
		totalPageStr = strings.TrimSpace(totalPageStr)
		totalPageStr = strings.TrimLeft(totalPageStr, "/ ")
		totalPageStr = strings.TrimRight(totalPageStr, "页")
		page.TotalPage, err = strconv.Atoi(totalPageStr)
	}
	return
}
