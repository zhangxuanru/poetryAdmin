/*
@Time : 2019/9/4 18:40
@Author : zxr
@File : catalogCategoryModel
@Software: GoLand
*/
package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

var TableBookContent = "poetry_ancient_book_content"

//poetry_ancient_book_content  古籍-正文内容
type BookContent struct {
	Id             int64  `orm:"column(id);auto"`
	BookId         int64  `orm:"column(book_id)"`
	CatalogId      int64  `orm:"column(catalog_id)"`
	Content        string `orm:"column(content)"`
	Translation    string `orm:"column(translation)"`
	TranslationId  int    `orm:"column(translation_id)"`
	TranslationUrl string `orm:"column(translation_url)"`
	AuthorId       int64  `orm:"column(author_id)"`
	SongUrl        string `orm:"column(song_url)"`
	SongFilePath   string `orm:"column(song_file_path)"`
	AddDate        int64  `orm:"column(add_date)"`
	UpdateDate     int64  `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(BookContent))
}

func (b *BookContent) TableName() string {
	return TableBookContent
}

func NewBookContent() *BookContent {
	return new(BookContent)
}

//保存正文数据
func (b *BookContent) Save(data *BookContent) (id int64, err error) {
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据目录ID查询正文内容
func (b *BookContent) GetContentByLogId(catalogId int64) (data BookContent, err error) {
	_, err = orm.NewOrm().QueryTable(TableBookContent).Filter("catalog_id", catalogId).All(&data, "id", "book_id", "author_id")
	return
}

//根据ID更新朗诵声音文件地址
func (b *BookContent) UpdateSongPath(contentId int64, filePath string) (id int64, err error) {
	data := &BookContent{
		Id:           contentId,
		SongFilePath: filePath,
		UpdateDate:   time.Now().Unix(),
	}
	id, err = orm.NewOrm().Update(data, "song_file_path", "update_date")
	return
}
