package models

import (
	"github.com/astaxie/beego/orm"
)

var TableAuthorData = "poetry_author_data"

//poetry_author_data 作者资料信息表
type AuthorData struct {
	Id         int   `orm:"column(id);auto"`
	AuthorId   int64 `orm:"column(author_id)"`
	DataId     int   `orm:"column(data_id)"`
	NotesId    int   `orm:"column(notes_id)"`
	Sort       int   `orm:"column(sort)"`
	AddDate    int64 `orm:"column(add_date)"`
	UpdateDate int64 `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(AuthorData))
}

func (a *AuthorData) TableName() string {
	return TableAuthorData
}

func NewAuthorData() *AuthorData {
	return new(AuthorData)
}

//保存作者资料信息
func (a *AuthorData) SaveAuthorData(data *AuthorData) (id int64, err error) {
	if data.Id > 0 {
		_, err = a.UpdateAuthorData(data)
		id = int64(data.Id)
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}

//更新作者资料信息
func (a *AuthorData) UpdateAuthorData(data *AuthorData, col ...string) (id int64, err error) {
	if len(col) == 0 {
		col = []string{"author_id", "data_id", "notes_id", "sort", "update_date"}
	}
	id, err = orm.NewOrm().Update(data, col...)
	return
}

//根据author_id|data_id查询资料数据
func (a *AuthorData) GetAuthorDataByDataId(authorId int64, dataId int64) (author AuthorData, err error) {
	_, err = orm.NewOrm().QueryTable(TableAuthorData).Filter("author_id", authorId).Filter("data_id", dataId).All(&author, "id", "notes_id", "sort")
	return
}
