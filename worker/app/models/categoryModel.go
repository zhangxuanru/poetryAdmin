package models

import (
	"github.com/astaxie/beego/orm"
	"poetryAdmin/worker/core/define"
	"time"
)

//poetry_category 诗文分类表
type Category struct {
	Id             int    `orm:"column(id);auto"`
	CatName        string `orm:"column(cat_name)"`
	SourceUrl      string `orm:"column(source_url)"`
	ShowPosition   int    `orm:"column(show_position)"`
	Pid            int    `orm:"column(pid)"`
	AddDate        int64  `orm:"column(add_date)"`
	LastUpdateTime int64  `orm:"column(last_update_time)"`
}

func init() {
	orm.RegisterModel(new(Category))
}

func (c *Category) TableName() string {
	return "poetry_category"
}

//保存分类
func InsertMultiCategoryByDataMap(data define.DataMap) (i int64, err error) {
	var categorys []Category
	for _, ret := range data {
		if ret.Text == "更多>>" {
			continue
		}
		category := Category{
			CatName:        ret.Text,
			SourceUrl:      ret.Href,
			ShowPosition:   int(ret.ShowPosition),
			AddDate:        time.Now().Unix(),
			LastUpdateTime: time.Now().Unix(),
		}
		categorys = append(categorys, category)
	}
	if len(categorys) > 0 {
		i, err = orm.NewOrm().InsertMulti(len(categorys), categorys)
	}
	return
}
