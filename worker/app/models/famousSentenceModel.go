/*
@Time : 2019/9/4 18:40
@Author : zxr
@File : catalogCategoryModel
@Software: GoLand
*/
package models

import (
	"github.com/astaxie/beego/orm"
)

var TableFamousSentence = "poetry_famous_sentence"

//poetry_famous_sentence  名句表
type FamousSentence struct {
	Id           int64  `orm:"column(id);auto"`
	CatId        int    `orm:"column(cat_id)"`
	Content      string `orm:"column(content)"`
	ContentCrc32 uint32 `orm:"column(content_crc32)"`
	PoetryTitle  string `orm:"column(poetry_title)"`
	PoetryId     int64  `orm:"column(poetry_id)"`
	AuthorId     int64  `orm:"column(author_id)"`
	Sort         int    `orm:"column(sort)"`
	SourceUrl    string `orm:"column(source_url)"`
	SourceCrc32  uint32 `orm:"column(source_crc32)"`
	AddDate      int64  `orm:"column(add_date)"`
	UpdateDate   int64  `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(FamousSentence))
}

func (f *FamousSentence) TableName() string {
	return TableFamousSentence
}

func NewFamousSentence() *FamousSentence {
	return new(FamousSentence)
}

//保存名句数据
func (f *FamousSentence) Save(data *FamousSentence) (id int64, err error) {
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据URL crc32值查询
func (f *FamousSentence) GetDataByCrc32(crc32 uint32, contentCrc uint32) (data FamousSentence, err error) {
	_, err = orm.NewOrm().QueryTable(TableFamousSentence).Filter("source_crc32", crc32).Filter("content_crc32", contentCrc).All(&data, "id")
	return
}
