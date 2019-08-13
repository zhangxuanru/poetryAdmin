package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

var TableCategoryGenre = "poetry_category_genre"

//poetry_category_genre 诗文类别体裁关联表
type CategoryGenre struct {
	Id         int   `orm:"column(id);auto"`
	CatId      int   `orm:"column(cat_id)"`
	GenreId    int64 `orm:"column(genre_id)"`
	UpdateDate int64 `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(CategoryGenre))
}

func NewCategoryGenre() *CategoryGenre {
	return new(CategoryGenre)
}

func (g *CategoryGenre) TableName() string {
	return TableCategoryGenre
}

//保存诗文类别体裁数据
func (g *CategoryGenre) SaveCategoryGenre(data *CategoryGenre) (id int64, err error) {
	var (
		genre CategoryGenre
	)
	data.UpdateDate = time.Now().Unix()
	genre, err = g.GetCateGenByIds(data.CatId, data.GenreId)
	if genre.Id > 0 {
		data.Id = genre.Id
		//_, _ = orm.NewOrm().Update(data, "update_date")
		return int64(genre.Id), nil
	}
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据分类ID和体裁ID查询数据
func (g *CategoryGenre) GetCateGenByIds(catId int, genreId int64) (data CategoryGenre, err error) {
	_, err = orm.NewOrm().QueryTable(TableCategoryGenre).Filter("cat_id", catId).Filter("genre_id", genreId).All(&data, "id", "update_date")
	return
}
