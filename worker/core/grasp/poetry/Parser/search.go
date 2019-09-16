/*
@Time : 2019/9/12 11:11
@Author : zxr
@File : search
@Software: GoLand
*/
package Parser

import (
	"github.com/PuerkitoBio/goquery"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"strings"
)

//解析搜索作者的HTML页面
//@see:https://so.gushiwen.org/search.aspx?value=%E7%AA%A6%E5%B7%A9
func ParserSearchAuthorHtml(html []byte) (author define.PoetryAuthorDetail, err error) {
	var query *goquery.Document
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	author.AuthorSrcUrl, _ = query.Find(".left>.sonspic>.cont>p").Eq(0).Find("a").Attr("href")
	author.AuthorContentUrl, _ = query.Find(".left>.sonspic>.cont>p").Eq(1).Find("a").Attr("href")
	if len(author.AuthorSrcUrl) > 0 && strings.Contains(author.AuthorSrcUrl, "http") == false {
		author.AuthorSrcUrl = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(author.AuthorSrcUrl, "/")
	}
	if len(author.AuthorContentUrl) > 0 && strings.Contains(author.AuthorContentUrl, "http") == false {
		author.AuthorContentUrl = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(author.AuthorContentUrl, "/")
	}
	return
}
