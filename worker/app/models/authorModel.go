package models

import (
	"github.com/astaxie/beego/orm"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"time"
)

//poetry_author作者表
type Author struct {
	Id           int    `orm:"column(id);auto"`
	Author       string `orm:"column(author)"`
	SourceUrl    string `orm:"column(source_url)"`
	DynastyId    int    `orm:"column(dynasty_id)"`
	AuthorsId    int    `orm:"column(authors_id)"`
	PhotoUrl     string `orm:"column(photo_url)"`
	PhotoId      int    `orm:"column(photo_id)"`
	AuthorDetail string `orm:"column(author_detail)"`
	PoetryCount  int    `orm:"column(poetry_count)"`
	IsRecommend  int    `orm:"column(is_recommend)"`
	Pinyin       string `orm:"column(pinyin)"`
	Acronym      string `orm:"column(acronym)"`
	AuthorTitle  string `orm:"column(author_title)"`
	AddDate      int64  `orm:"column(add_date)"`
	UpdateDate   int64  `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(Author))
}

func (c *Author) TableName() string {
	return "poetry_author"
}

//根据首页的数据保存作者信息
func InsertMultiAuthorByDataMap(data define.DataMap) (i int64, err error) {
	var authors []Author
	var acronym string
	for _, ret := range data {
		if ret.Text == "更多>>" {
			continue
		}
		pinyin := tools.PinYin(ret.Text)
		if len(pinyin) > 0 {
			acronym = pinyin[:1]
		}
		au := Author{
			Author:      ret.Text,
			SourceUrl:   ret.Href,
			IsRecommend: 1,
			AddDate:     time.Now().Unix(),
			UpdateDate:  time.Now().Unix(),
			Pinyin:      pinyin,
			Acronym:     acronym,
		}
		authors = append(authors, au)
	}
	if len(authors) > 0 {
		i, err = orm.NewOrm().InsertMulti(len(authors), authors)
	}
	return
}
