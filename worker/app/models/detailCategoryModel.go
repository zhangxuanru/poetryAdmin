package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
)

var TableDetailCategory = "poetry_detail_category"

//poetry_detail_category 诗文详情分类表
type DetailCategory struct {
	Id         int   `orm:"column(id);auto"`
	PoetryId   int   `orm:"column(poetry_id)"`
	CategoryId int   `orm:"column(category_id)"`
	UpdateTime int64 `orm:"column(update_time)"`
}

func init() {
	orm.RegisterModel(new(DetailCategory))
}

func (d *DetailCategory) TableName() string {
	return TableDetailCategory
}

func NewDetailCategory() *DetailCategory {
	return new(DetailCategory)
}

//保存诗词分类信息
func (d *DetailCategory) SaveDetailCategory(data *DetailCategory) (id int64, err error) {
	var (
		category DetailCategory
	)
	if data.PoetryId == 0 || data.CategoryId == 0 {
		return 0, errors.New("PoetryId or CategoryId eq 0")
	}
	if category, err = d.FindDataByCatId(data.CategoryId, data.PoetryId); err != nil {
		return 0, err
	}
	if category.Id > 0 {
		return int64(category.Id), nil
	}
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据分类ID和诗词ID查询是否已存在
func (d *DetailCategory) FindDataByCatId(categoryId, poetryId int) (data DetailCategory, err error) {
	_, err = orm.NewOrm().QueryTable(TableDetailCategory).Filter("poetry_id", poetryId).Filter("category_id", categoryId).All(&data, "id")
	return
}
