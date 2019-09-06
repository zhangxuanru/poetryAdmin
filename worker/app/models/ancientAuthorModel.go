package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

var TableAncientAuthor = "poetry_ancient_author"

//poetry_ancient_author 古籍作者表
type AncientAuthor struct {
	Id         int    `orm:"column(id);auto"`
	AuthorName string `orm:"column(author_name)"`
	SourceUrl  string `orm:"column(source_url)"`
	AddDate    int64  `orm:"column(add_date)"`
}

func init() {
	orm.RegisterModel(new(AncientAuthor))
}

func (a *AncientAuthor) TableName() string {
	return TableAncientAuthor
}

func NewAncientAuthor() *AncientAuthor {
	return new(AncientAuthor)
}

//根据作者姓名查询作者信息
func (a *AncientAuthor) GetAuthorDataByAuthorName(authorName string) (author AncientAuthor, err error) {
	_, err = orm.NewOrm().QueryTable(TableAncientAuthor).Filter("author_name", authorName).All(&author, "id")
	return
}

//保存作者信息
func (a *AncientAuthor) SaveAuthor(data *AncientAuthor) (id int64, err error) {
	var (
		author AncientAuthor
	)
	if data.AuthorName == "" {
		return 0, nil
	}
	if author, err = a.GetAuthorDataByAuthorName(data.AuthorName); err != nil {
		return 0, err
	}
	if author.Id > 0 {
		return int64(author.Id), nil
	}
	data.AddDate = time.Now().Unix()
	id, err = orm.NewOrm().Insert(data)
	return
}
