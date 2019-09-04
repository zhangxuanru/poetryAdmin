/*
@Time : 2019/9/4 18:40
@Author : zxr
@File : catalogCategoryModel
@Software: GoLand
*/
package models

import "github.com/astaxie/beego/orm"

var TableBookCatalogueCategory = "poetry_ancient_book_catalog_category"

//poetry_ancient_book_catalog_category  古籍-书名目录分类表
type BookCatalogueCategory struct {
	Id      int64  `orm:"column(id);auto"`
	BookId  int64  `orm:"column(book_id)"`
	CatName string `orm:"column(cat_name)"`
	Sort    int    `orm:"column(sort)"`
	AddDate int64  `orm:"column(add_date)"`
}

func init() {
	orm.RegisterModel(new(BookCatalogueCategory))
}

func (b *BookCatalogueCategory) TableName() string {
	return TableBookCatalogueCategory
}

func NewBookCatalogueCategory() *BookCatalogueCategory {
	return new(BookCatalogueCategory)
}

//保存目录分类数据
func (b *BookCatalogueCategory) Save(data *BookCatalogueCategory) (id int64, err error) {
	var (
		catalogue BookCatalogueCategory
	)
	if catalogue, err = b.GetCatalogueCategory(data.BookId, data.CatName); err != nil {
		return 0, err
	}
	if catalogue.Id > 0 {
		return catalogue.Id, nil
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}

//根据书ID和分类名查询是否有数据
func (b *BookCatalogueCategory) GetCatalogueCategory(bookId int64, catName string) (data BookCatalogueCategory, err error) {
	_, err = orm.NewOrm().QueryTable(TableBookCatalogueCategory).Filter("book_id", bookId).Filter("cat_name", catName).All(&data, "id")
	return
}
