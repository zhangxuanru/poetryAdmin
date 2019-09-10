/*
@Time : 2019/9/6 18:56
@Author : zxr
@File : famous
@Software: GoLand
*/
package define

//名句--主题
type Classify struct {
	ThemeTitle  string
	Title       string
	LinkUrl     string
	Sort        int
	ContentList []Content
}

//主题下的分类信息 一个主题下有多个分类
type ThemeCategory struct {
	Title        string
	LinkUrl      string
	ClassifyList []*Classify
}

//名句-正文内容
type Content struct {
	Text        string
	LinkUrl     string
	Sort        int
	PoetryText  string
	PoetryTitle string
	PoetryLink  string
	AuthorName  string
}

//名句-正文列表分页
type ContentPage struct {
	NextPageUrl string
	TotalPage   int
	IsNextPage  bool
}
