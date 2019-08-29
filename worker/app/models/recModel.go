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

var TableRec = "poetry_content_apprec"

//poetry_content_apprec 诗词详情赏析信息关联表
type ContentRec struct {
	Id         int   `orm:"column(id);auto"`
	PoetryId   int   `orm:"column(poetry_id)"`
	ApprecId   int   `orm:"column(apprec_id)"`
	NotesId    int64 `orm:"column(notes_id)"`
	Sort       int   `orm:"column(sort)"`
	AddDate    int64 `orm:"column(add_date)"`
	UpdateDate int64 `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(ContentRec))
}

func (c *ContentRec) TableName() string {
	return TableRec
}

func NewContentRec() *ContentRec {
	return new(ContentRec)
}

//保存赏析信息,如果存在则直接返回
func (c *ContentRec) SaveResData(data *ContentRec) (id int64, err error) {
	var rec ContentRec
	if data.PoetryId == 0 || data.ApprecId == 0 {
		return 0, errors.New("PoetryId or ApprecId eq 0")
	}
	if rec, err = c.FindNotesId(data.PoetryId, data.ApprecId); err != nil {
		return 0, err
	}
	if rec.Id > 0 {
		return int64(rec.Id), nil
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}

//保存赏析信息,直接插入，不管是否存在
func (c *ContentRec) InsertRec(data *ContentRec) (id int64, err error) {
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据诗词ID和赏析ID查询文本ID
func (c *ContentRec) FindNotesId(poetryId, recId int) (data ContentRec, err error) {
	_, err = orm.NewOrm().QueryTable(TableRec).Filter("poetry_id", poetryId).Filter("apprec_id", recId).All(&data, "id", "notes_id", "sort")
	return
}
