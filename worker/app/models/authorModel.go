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
	Id            int    `orm:"column(id);auto"`
	Author        string `orm:"column(author)"`
	SourceUrl     string `orm:"column(source_url)"`
	WorksUrl      string `orm:"column(works_url)"`
	DynastyId     int    `orm:"column(dynasty_id)"`
	AuthorsId     int    `orm:"column(authors_id)"`
	PhotoUrl      string `orm:"column(photo_url)"`
	PhotoFileName string `orm:"column(photo_file_name)"`
	AuthorIntro   string `orm:"column(author_intro)"`
	PoetryCount   int    `orm:"column(poetry_count)"`
	IsRecommend   int    `orm:"column(is_recommend)"`
	Pinyin        string `orm:"column(pinyin)"`
	Acronym       string `orm:"column(acronym)"`
	AuthorTitle   string `orm:"column(author_title)"`
	AddDate       int64  `orm:"column(add_date)"`
	UpdateDate    int64  `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(Author))
}

func (c *Author) TableName() string {
	return TableAuthor
}

func NewAuthor() *Author {
	return new(Author)
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
	_, err = orm.NewOrm().QueryTable(TableAuthor).Filter("author", authorName).All(&author, "id", "pinyin", "dynasty_id")
	return
}

//根据作者姓名和朝代查询作者信息
func GetAuthorByNameAndDyId(authorName string, dynastyId int) (author Author, err error) {
	_, err = orm.NewOrm().QueryTable(TableAuthor).Filter("author", authorName).Filter("dynasty_id", dynastyId).All(&author, "id", "pinyin", "dynasty_id")
	return
}

//保存作者信息
func (a *Author) SaveAuthor(data *Author) (id int64, err error) {
	var (
		author  Author
		acronym string
	)
	if data.Author == "" {
		return 0, nil
	}
	if author, err = GetAuthorDataByAuthorName(data.Author); err != nil {
		return 0, err
	}
	if author.Id > 0 && author.DynastyId > 0 {
		return int64(author.Id), nil
	}
	pinyin := tools.PinYin(data.Author)
	if len(pinyin) > 0 {
		acronym = pinyin[:1]
	}
	data.AddDate = time.Now().Unix()
	data.UpdateDate = time.Now().Unix()
	data.Pinyin = pinyin
	data.Acronym = acronym
	if author.Id > 0 {
		data.Id = author.Id
		if data.DynastyId > 0 {
			_, _ = orm.NewOrm().Update(data, "dynasty_id")
		}
		if len(data.SourceUrl) > 0 {
			_, _ = orm.NewOrm().Update(data, "source_url")
		}
		if len(data.WorksUrl) > 0 {
			_, _ = orm.NewOrm().Update(data, "works_url")
		}
		if data.PoetryCount > 0 {
			_, _ = orm.NewOrm().Update(data, "poetry_count")
		}
		if len(data.PhotoUrl) > 0 {
			_, _ = orm.NewOrm().Update(data, "photo_url")
		}
		return int64(author.Id), nil
	}
	id, err = orm.NewOrm().Insert(data)
	return
}

//更新作者信息
func (a *Author) UpdateAuthor(data *Author, fields ...string) (id int64, err error) {
	var (
		author  Author
		acronym string
		pinyin  string
	)
	if author, err = GetAuthorDataByAuthorName(data.Author); err != nil {
		return
	}
	if len(fields) == 0 {
		fields = []string{"update_date", "source_url", "works_url", "dynasty_id", "photo_url", "photo_file_name", "author_intro", "poetry_count"}
		if len(author.Pinyin) == 0 {
			fields = append(fields, "pinyin", "acronym")
		}
	}
	if pinyin = tools.PinYin(data.Author); len(pinyin) > 0 {
		acronym = pinyin[:1]
		data.Pinyin = pinyin
		data.Acronym = acronym
	}
	data.AddDate = time.Now().Unix()
	data.UpdateDate = time.Now().Unix()
	if author.Id > 0 {
		data.Id = author.Id
		id = int64(data.Id)
		_, err = orm.NewOrm().Update(data, fields...)
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	if err != nil {
		return
	}
	if id > 0 {
		return id, nil
	}
	return 0, err
}

//更新头像路径
func (a *Author) UpdateAuthorPhoto(data *Author) (id int64, err error) {
	var (
		author Author
	)
	if author, err = GetAuthorDataByAuthorName(data.Author); err != nil || author.Id == 0 {
		return
	}
	data.Id = author.Id
	id, err = orm.NewOrm().Update(data, "photo_file_name")
	return
}

func (a *Author) UpdateAuthorDynasty(data *Author) (int64, error) {
	return orm.NewOrm().Update(data, "dynasty_id")
}

func (a *Author) GetOrm() orm.Ormer {
	return orm.NewOrm()
}
