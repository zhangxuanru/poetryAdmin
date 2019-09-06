/*
@Time : 2019/9/6 18:56
@Author : zxr
@File : famous
@Software: GoLand
*/
package define

//名句--主题分类
type Classify struct {
	Title       string
	LinkUrl     string
	Sort        int
	ContentList []Content
}

//名句-正文内容
type Content struct {
	Text       string
	LinkUrl    string
	Sort       int
	PoetryText string
	AuthorName string
}
