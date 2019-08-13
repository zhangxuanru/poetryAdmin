package models

import "github.com/astaxie/beego/orm"

var TableGenre = "poetry_genre"

//poetry_genre诗文体裁表
type Genre struct {
	Id        int    `orm:"column(id);auto"`
	GenreName string `orm:"column(genre_name)"`
	AddDate   int64  `orm:"column(add_date)"`
}

func init() {
	orm.RegisterModel(new(Genre))
}

func (g *Genre) TableName() string {
	return TableGenre
}

//保存诗文体裁
func SaveGenre(data *Genre) (id int64, err error) {
	var (
		genreData Genre
	)
	genreData, err = GetGemreByName(data.GenreName)
	if genreData.Id > 0 {
		return int64(genreData.Id), nil
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}

//根据名字搜索体裁是否存在
func GetGemreByName(name string) (genreData Genre, err error) {
	_, err = orm.NewOrm().QueryTable(TableGenre).Filter("genre_name", name).All(&genreData, "id")
	return
}
