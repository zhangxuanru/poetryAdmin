/*
@Time : 2019/8/28 19:22
@Author : zxr
@File : transModel
@Software: GoLand
*/
package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
)

var TableTrans = "poetry_content_trans"

//poetry_content_trans 诗词详情翻译信息关联表
type ContentTrans struct {
	Id         int   `orm:"column(id);auto"`
	PoetryId   int   `orm:"column(poetry_id)"`
	TransId    int   `orm:"column(trans_id)"`
	NotesId    int64 `orm:"column(notes_id)"`
	Sort       int   `orm:"column(sort)"`
	AddDate    int64 `orm:"column(add_date)"`
	UpdateDate int64 `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(ContentTrans))
}

func (c *ContentTrans) TableName() string {
	return TableTrans
}

func NewContentTrans() *ContentTrans {
	return new(ContentTrans)
}

//保存翻译信息,如果存在则直接返回
func (c *ContentTrans) SaveTrans(data *ContentTrans) (id int64, err error) {
	var trans ContentTrans
	if data.PoetryId == 0 || data.TransId == 0 {
		return 0, errors.New("PoetryId or TransId eq 0")
	}
	if trans, err = c.FindNotesId(data.PoetryId, data.TransId); err != nil {
		return 0, err
	}
	if trans.Id > 0 {
		return int64(trans.Id), nil
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}

//保存翻译信息,直接插入，不管是否存在
func (c *ContentTrans) InsertTrans(data *ContentTrans) (id int64, err error) {
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据诗词ID和翻译ID查询文本ID
func (c *ContentTrans) FindNotesId(poetryId, transId int) (data ContentTrans, err error) {
	_, err = orm.NewOrm().QueryTable(TableTrans).Filter("poetry_id", poetryId).Filter("trans_id", transId).All(&data, "id", "notes_id", "sort")
	return
}

func (c *ContentTrans) GetOrm() orm.Ormer {
	return orm.NewOrm()
}
