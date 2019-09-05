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
	GuWenCategory
}

//书籍目录信息
type CataLogData struct {
	CateName    string
	Sort        int
	CatalogList []CataLog
}

type CataLog struct {
	CatalogTitle string
	Sort         int
	LinkUrl      string
	BookId       int64
}

//书籍信息和对应目录
type BookCatalogue struct {
	BookTitle   string
	BookLinkUrl string
	CatalogList []CataLog
}

//书籍目录详情
type BookCatalogueContent struct {
	ShortCatalogueTitle string //目录短标题
	CatalogueLinkUrl    string //目录链接
	CatalogueTitle      string //目录长标题
	SongId              int    //声音ID
	SongUrl             string //声音URL文件
	AuthorName          string //作者姓名
	AuthorLinkUrl       string //作者链接地址
	Content             string //正文内容
	Translation         string //译文内容
	TranslationUrl      string //译文URL
	TranslationId       int    //译文ID
}
