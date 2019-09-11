/*
@Time : 2019/9/5 18:05
@Author : zxr
@File : content
@Software: GoLand
*/
package Parser

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"strconv"
	"strings"
)

//解析目录内容详情
func ParseGuWenContent(html []byte) (content *define.BookCatalogueContent, err error) {
	var (
		query  *goquery.Document
		bookId string
		data   define.BookCatalogueContent
	)
	if len(html) == 0 {
		return content, errors.New("html is nil")
	}
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	data.CatalogueTitle = query.Find(".sons>.cont>h1>span>b").Text()
	if bookIdStr, ok := query.Find(".sons>.cont>h1>a").Eq(0).Attr("href"); ok {
		if bookId = tools.TrimPlayBookKv(bookIdStr); len(bookId) > 0 {
			data.SongId, err = strconv.Atoi(bookId)
		}
	}
	if data.SongId > 0 {
		data.SongUrl = config.G_Conf.GushiwenAncientUrl + "bookvplay.aspx?id=" + bookId
	}
	if href, ok := query.Find(".sons>.cont>h1>a").Eq(1).Attr("href"); ok {
		if yiZhuIdStr := tools.TrimYiZhu(href); len(yiZhuIdStr) > 0 {
			data.TranslationId, _ = strconv.Atoi(yiZhuIdStr)
			data.TranslationUrl = config.G_Conf.GushiwenAncientUrl + "ajaxbfanyi.aspx?id=" + yiZhuIdStr
		}
	}
	data.AuthorName = query.Find(".sons>.cont>.source>a").Text()
	data.AuthorLinkUrl, _ = query.Find(".sons>.cont>.source>a").Attr("href")
	data.Content, err = query.Find(".sons>.cont>.contson").Html()
	data.Content = strings.TrimSpace(data.Content)
	//data.Content = tools.TrimHtml(data.Content)
	return &data, err
}
