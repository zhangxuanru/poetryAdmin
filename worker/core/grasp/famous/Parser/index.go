/*
@Time : 2019/9/6 18:53
@Author : zxr
@File : index
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

//解析名句首页数据格式
//https://so.gushiwen.org/mingju/
func ParseFamousIndexClassify(html []byte) (classifyList []define.Classify, err error) {
	var (
		query    *goquery.Document
		classify define.Classify
	)
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	query.Find(".left>.titletype>.son2>.sright>a").Each(func(i int, selection *goquery.Selection) {
		classify.Title = strings.TrimSpace(selection.Text())
		classify.LinkUrl, _ = selection.Attr("href")
		classify.Sort = i
		if len(classify.LinkUrl) > 0 {
			classify.LinkUrl = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(classify.LinkUrl, "/")
		}
		classifyList = append(classifyList, classify)
	})
	return classifyList, nil
}
