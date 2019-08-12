package models

import (
	"github.com/astaxie/beego/orm"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"time"
)

var TableAuthor = "poetry_author"

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
	return TableAuthor
}

//根据首页的数据保存作者信息
func InsertMultiAuthorByDataMap(data define.DataMap) (i int64, err error) {
	var (
		authors   []Author
		authorRow Author
		acronym   string
	)
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
		authorRow, _ = GetAuthorDataByAuthorName(au.Author)
		if authorRow.Id > 0 {
			au.Id = authorRow.Id
			_, _ = orm.NewOrm().Update(&au, "source_url", "update_date", "pinyin", "acronym")
			au.Id = 0
		} else {
			authors = append(authors, au)
		}
	}
	if len(authors) > 0 {
		i, err = orm.NewOrm().InsertMulti(len(authors), authors)
	}
	return
}

//根据作者姓名查询作者信息
func GetAuthorDataByAuthorName(authorName string) (author Author, err error) {
	_, err = orm.NewOrm().QueryTable(TableAuthor).Filter("author", authorName).All(&author)
	return
}
