package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
	"time"
)

var TableDynasty = "poetry_dynasty"

//poetry_dynasty朝代表
type Dynasty struct {
	Id          int    `orm:"column(id);auto"`
	DynastyName string `orm:"column(dynasty_name)"`
	AddDate     int64  `orm:"column(add_date)"`
}

func init() {
	orm.RegisterModel(new(Dynasty))
}

func (d *Dynasty) TableName() string {
	return TableDynasty
}

func NewDynasty() *Dynasty {
	return new(Dynasty)
}

//根据朝代名查询朝代ID
func (d *Dynasty) GetOneByName(name string) (data Dynasty, err error) {
	if len(name) == 0 {
		return
	}
	_, err = orm.NewOrm().QueryTable(TableDynasty).Filter("dynasty_name", name).All(&data, "id")
	return
}

//插入数据
func (d *Dynasty) SaveName(name string) (id int64, err error) {
	data := &Dynasty{
		DynastyName: name,
		AddDate:     time.Now().Unix(),
	}
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据名字查询是否保存，如果没有保存则 插入，保存了则返回ID
func (d *Dynasty) GetIdBySaveName(name string) (id int64, err error) {
	var data Dynasty
	if len(name) == 0 {
		return 0, errors.New("name is nil")
	}
	if data, err = d.GetOneByName(name); err != nil {
		return
	}
	if data.Id > 0 {
		return int64(data.Id), nil
	}
	return d.SaveName(name)
}
