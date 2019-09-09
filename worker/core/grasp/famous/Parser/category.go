/*
@Time : 2019/9/9 15:08
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
	"strings"
)

//解析名句-主题下的分类信息
func ParseFamousThemeClassify(html []byte) (classifyList []*define.Classify, err error) {
	var query *goquery.Document
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	query.Find(".left>.titletype>.son2").Eq(1).Find(".sright>a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		if len(href) > 0 {
			href = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(href, "/")
		}
		classify := &define.Classify{
			Title:   strings.TrimSpace(selection.Text()),
			LinkUrl: strings.TrimSpace(href),
			Sort:    i,
		}
		classifyList = append(classifyList, classify)
	})
	return
}
