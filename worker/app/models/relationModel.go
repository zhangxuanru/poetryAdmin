package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

var TableContentRelation = "poetry_content_relation"

//poetry_content_relation 诗词关联表
type ContentRelation struct {
	Id         int   `orm:"column(id);auto"`
	PoetryId   int64 `orm:"column(poetry_id)"`
	AuthorId   int64 `orm:"column(author_id)"`
	CategoryId int64 `orm:"column(category_id)"`
	GenreId    int64 `orm:"column(genre_id)"`
	Form       int   `orm:"column(form)"`
	Sort       int   `orm:"column(sort)"`
	AddDate    int64 `orm:"column(add_date)"`
	UpdateDate int64 `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(ContentRelation))
}

func NewContentRelation() *ContentRelation {
	return new(ContentRelation)
}

func (c *ContentRelation) TableName() string {
	return TableContentRelation
}

//保存诗词关联关系
func (c *ContentRelation) SaveContentRelation(data *ContentRelation) (id int64, err error) {
	var content ContentRelation
	if content, err = c.GetDataByMoreId(data.PoetryId, data.CategoryId, data.GenreId, data.AuthorId); err != nil {
		return 0, err
	}
	data.UpdateDate = time.Now().Unix()
	data.AddDate = time.Now().Unix()
	if content.Id > 0 {
		data.Id = content.Id
		_, err = c.UpdateRelation(data)
		return int64(content.Id), err
	}
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据ID查询关联数据
func (c *ContentRelation) GetDataByMoreId(poetryId, categoryId, genreId, authorId int64) (content ContentRelation, err error) {
	_, err = orm.NewOrm().QueryTable(TableContentRelation).Filter("poetry_id", poetryId).Filter("category_id", categoryId).Filter("genre_id", genreId).Filter("author_id", authorId).All(&content, "id", "update_date")
	return
}

//更新关系
func (c *ContentRelation) UpdateRelation(data *ContentRelation, col ...string) (id int64, err error) {
	if len(col) == 0 {
		col = []string{"poetry_id", "author_id", "category_id", "genre_id", "form", "sort", "update_date"}
	}
	id, err = orm.NewOrm().Update(data, col...)
	return
}
