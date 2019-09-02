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
	"poetryAdmin/worker/core/define"
	"strings"
)

//解析古文首页数据格式
// https://so.gushiwen.org/guwen/
func ParseGuWenIndexCategory(html []byte) (categoryData []define.GuWenCategoryList, err error) {
	var (
		query        *goquery.Document
		categoryList define.GuWenCategoryList
		nodeCategory []define.GuWenCategory
	)
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	//抓取分类
	query.Find(".titletype>.son2").Each(func(i int, selection *goquery.Selection) {
		parentCateText := selection.Find(".sleft>a").Text()
		link, _ := selection.Find(".sleft>a").Attr("href")
		categoryList.CategoryName = strings.TrimSpace(parentCateText)
		categoryList.LinkUrl = strings.TrimSpace(link)
		categoryList.Sort = i
		selection.Find(".sright>a").Each(func(i int, selection *goquery.Selection) {
			href, _ := selection.Attr("href")
			nodeCategory = append(nodeCategory, define.GuWenCategory{
				CategoryName: strings.TrimSpace(selection.Text()),
				LinkUrl:      strings.TrimSpace(href),
				Sort:         i,
			})
		})
		categoryList.SubNode = nodeCategory
		categoryData = append(categoryData, categoryList)
		nodeCategory = nil
	})
	return categoryData, nil
}
