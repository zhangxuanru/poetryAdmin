/*
@Time : 2019/8/30 18:17
@Author : zxr
@File : index
@Software: GoLand
*/
package Parser

import (
	"github.com/PuerkitoBio/goquery"
	"poetryAdmin/worker/app/tools"
)

//解析首页数据格式
// https://so.gushiwen.org/guwen/
func ParseIndex(html []byte) (err error) {
	var query *goquery.Document
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	//抓取分类
	return
}
