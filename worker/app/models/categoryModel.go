package models

import (
	"github.com/astaxie/beego/orm"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"time"
)

var TableCategory = "poetry_category"

type uintMaps map[uint32]Category

//poetry_category 诗文分类表
type Category struct {
	Id             int    `orm:"column(id);auto"`
	CatName        string `orm:"column(cat_name)"`
	SourceUrl      string `orm:"column(source_url)"`
	SourceUrlCrc32 uint32 `orm:"column(source_url_crc32)"`
	ShowPosition   int    `orm:"column(show_position)"`
	Pid            int    `orm:"column(pid)"`
	AddDate        int64  `orm:"column(add_date)"`
	LastUpdateTime int64  `orm:"column(last_update_time)"`
}

func init() {
	orm.RegisterModel(new(Category))
}

func (c *Category) TableName() string {
	return TableCategory
}

//保存更新分类
func InsertMultiCategoryByDataMap(data define.DataMap) (i int64, err error) {
	var (
		categorys []Category
	)
	for _, ret := range data {
		if ret.Text == "更多>>" {
			continue
		}
		category := Category{
			CatName:        ret.Text,
			SourceUrl:      ret.Href,
			SourceUrlCrc32: tools.Crc32(ret.Href),
			ShowPosition:   int(ret.ShowPosition),
			AddDate:        time.Now().Unix(),
			LastUpdateTime: time.Now().Unix(),
		}
		categories, _ := GetDataByCateName(category.CatName, category.ShowPosition)
		if categories.Id > 0 {
			category.Id = categories.Id
			_, _ = orm.NewOrm().Update(&category, "cat_name", "source_url", "source_url_crc32", "show_position", "pid", "last_update_time")
			category.Id = 0
		} else {
			categorys = append(categorys, category)
		}
	}
	if len(categorys) > 0 {
		i, err = orm.NewOrm().InsertMulti(len(categorys), categorys)
	}
	return
}

//根据位置查询分类数据
func GetCategoryDataByPosition(showPosition define.ShowPosition) (maps uintMaps, err error) {
	var categorys []Category
	maps = make(uintMaps)
	_, err = orm.NewOrm().QueryTable(TableCategory).Filter("show_position", showPosition).All(&categorys, "id", "cat_name", "source_url", "pid", "source_url_crc32")
	if len(categorys) > 0 {
		for _, val := range categorys {
			maps[val.SourceUrlCrc32] = val
		}
	}
	return
}

//根据分类名和位置查询数据
func GetDataByCateName(cateName string, position int) (categorys Category, err error) {
	_, err = orm.NewOrm().QueryTable(TableCategory).Filter("show_position", position).Filter("cat_name", cateName).All(&categorys, "id")
	return categorys, err
}

//根据分类名和位置和PID查询数据
func GetDataByCateNameAndPid(cateName string, position int, pid int) (category Category, err error) {
	_, err = orm.NewOrm().QueryTable(TableCategory).Filter("show_position", position).Filter("cat_name", cateName).Filter("pid", pid).All(&category, "id")
	return category, err
}

//保存分类数据
func SaveCategoryData(data *Category) (id int64, err error) {
	id, err = orm.NewOrm().Insert(data)
	return
}
