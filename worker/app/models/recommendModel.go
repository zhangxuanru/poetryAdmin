/*
@Time : 2019/9/16 17:07
@Author : zxr
@File : recommendModel
@Software: GoLand
*/
package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
)

var TableRecommend = "poetry_recommend"

//poetry_recommend 诗词推荐表
type Recommend struct {
	Id          int64 `orm:"column(id);auto"`
	PoetryId    int64 `orm:"column(poetry_id)"`
	Sort        int   `orm:"column(sort)"`
	Status      int   `orm:"column(status)"`
	RecommeTime int64 `orm:"column(recomme_time)"`
	AddDate     int64 `orm:"column(add_date)"`
}

func init() {
	orm.RegisterModel(new(Recommend))
}

func (r *Recommend) TableName() string {
	return TableRecommend
}

func NewRecommend() *Recommend {
	return new(Recommend)
}

func (r *Recommend) SaveRecommend(data *Recommend) (id int64, err error) {
	var recommend Recommend
	if data.PoetryId == 0 {
		return 0, errors.New("PoetryId is nil")
	}
	if recommend, err = r.FindIdByPoetryIdAndTime(data.PoetryId, data.RecommeTime); err != nil {
		return
	}
	if recommend.Id > 0 {
		return recommend.Id, nil
	}
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据诗词ID和推荐时间查询是否存在
func (r *Recommend) FindIdByPoetryIdAndTime(poetryId int64, recommendTime int64) (data Recommend, err error) {
	_, err = orm.NewOrm().QueryTable(TableRecommend).Filter("poetry_id", poetryId).Filter("recomme_time", recommendTime).All(&data, "id")
	return
}
