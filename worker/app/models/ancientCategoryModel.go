/*
@Time : 2019/9/2 15:51
@Author : zxr
@File : ancientCategoryModel
@Software: GoLand
*/
package models

import "github.com/astaxie/beego/orm"

var TableAncientCategory = "poetry_ancient_category"

//poetry_ancient_category 古籍-栏目分类表
type AncientCategory struct {
	Id      int    `orm:"column(id);auto"`
	CatName string `orm:"column(cat_name)"`
	SrcUrl  string `orm:"column(src_url)"`
	Pid     int64  `orm:"column(pid)"`
	Sort    int    `orm:"column(sort)"`
	AddDate int64  `orm:"column(add_date)"`
}

func init() {
	orm.RegisterModel(new(AncientCategory))
}

func (a *AncientCategory) TableName() string {
	return TableAncientCategory
}

func NewAncientCategory() *AncientCategory {
	return new(AncientCategory)
}

//写入分类表
func (a *AncientCategory) InsertCategory(data *AncientCategory) (id int64, err error) {
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据分类名称查询数据
func (a *AncientCategory) GetCategoryDataByName(catName string) (data AncientCategory, err error) {
	_, err = orm.NewOrm().QueryTable(TableAncientCategory).Filter("cat_name", catName).All(&data, "id", "pid", "sort")
	return
}

//保存分类，有就更新，没有就插入
func (a *AncientCategory) SaveCategory(data *AncientCategory) (id int64, err error) {
	var (
		category AncientCategory
	)
	if category, err = a.GetCategoryDataByName(data.CatName); err != nil {
		return
	}
	if category.Id > 0 {
		//data.Id = category.Id
		//_, err = orm.NewOrm().Update(data, "cat_name", "src_url", "pid", "sort")
		return int64(category.Id), err
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}
