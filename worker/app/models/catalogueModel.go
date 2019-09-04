/*
@Time : 2019/9/4 18:40
@Author : zxr
@File : catalogCategoryModel
@Software: GoLand
*/
package models

import "github.com/astaxie/beego/orm"

var TableBookCatalogue = "poetry_ancient_book_catalog"

//poetry_ancient_book_catalog  古籍-书名目录表
type BookCatalogue struct {
	Id               int64  `orm:"column(id);auto"`
	BookId           int64  `orm:"column(book_id)"`
	CatalogTitle     string `orm:"column(catalog_title)"`
	CatalogCatgoryId int64  `orm:"column(catalog_catgory_id)"`
	LinkUrl          string `orm:"column(link_url)"`
	Sort             int    `orm:"column(sort)"`
	AddDate          int64  `orm:"column(add_date)"`
}

func init() {
	orm.RegisterModel(new(BookCatalogue))
}

func (b *BookCatalogue) TableName() string {
	return TableBookCatalogue
}

func NewBookCatalogue() *BookCatalogue {
	return new(BookCatalogue)
}

//保存目录数据
func (b *BookCatalogue) Save(data *BookCatalogue) (id int64, err error) {
	var (
		catalogue BookCatalogue
	)
	if catalogue, err = b.GetCatalogueData(data.BookId, data.CatalogTitle); err != nil {
		return 0, err
	}
	if catalogue.Id > 0 {
		return catalogue.Id, nil
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}

//根据书ID和目录标题查询是否有数据
func (b *BookCatalogue) GetCatalogueData(bookId int64, title string) (data BookCatalogue, err error) {
	_, err = orm.NewOrm().QueryTable(TableBookCatalogue).Filter("book_id", bookId).Filter("catalog_title", title).All(&data, "id")
	return
}
