/*
@Time : 2019/9/2 15:00
@Author : zxr
@File : ancient
@Software: GoLand
*/
package define

//古文首页分类结构体
type GuWenCategoryList struct {
	GuWenCategory
	SubNode []GuWenCategory //子分类
}

type GuWenCategory struct {
	CategoryName string
	LinkUrl      string
	Sort         int
}

//分类下的书籍信息HTML
type GuWenCategoryBookHtml struct {
	GuWenCategory
	Html []byte
}

//分类下的书籍信息
type GuWenBookCover struct {
	Title        string
	LinkUrl      string
	SongUrl      string
	SongSrcUrl   string
	Introduction string
	FamousTotal  int
	CoverChart   string
}
