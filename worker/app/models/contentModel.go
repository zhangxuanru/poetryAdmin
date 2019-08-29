package models

import (
	"errors"
	"github.com/astaxie/beego/orm"
)

var TableContent = "poetry_content"

//poetry_content 诗词表
type Content struct {
	Id             int    `orm:"column(id);auto"`
	Title          string `orm:"column(title)"`
	Content        string `orm:"column(content)"`
	AuthorId       int64  `orm:"column(author_id)"`
	SourceUrl      string `orm:"column(source_url)"`
	SourceUrlCrc32 uint32 `orm:"column(sourceurl_crc32)"`
	GenreId        int64  `orm:"column(genre_id)"`
	CreatBackId    int64  `orm:"column(creat_back_id)"`
	Sort           int    `orm:"column(sort)"`
	AddDate        int64  `orm:"column(add_date)"`
	UpdateDate     int64  `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(Content))
}

func NewContent() *Content {
	return new(Content)
}

func (c *Content) TableName() string {
	return TableContent
}

//保存诗词内容
func (c *Content) SaveContent(data *Content) (id int64, err error) {
	var content Content
	if len(data.Title) == 0 {
		return 0, errors.New("title is nil")
	}
	if data.SourceUrlCrc32 > 0 {
		if content, err = c.GetContentByCrc32(data.SourceUrlCrc32); err != nil {
			return 0, err
		}
	}
	//if content.Id == 0 {
	//	if content, err = c.GetByTitleAuthorId(data.Title, data.AuthorId); err != nil {
	//		return 0, err
	//	}
	//}
	if content.Id > 0 {
		if len(data.Content) > 0 {
			data.Id = content.Id
			_, err = c.UpdateContent(data, "title", "content", "source_url", "sourceurl_crc32", "author_id", "genre_id", "creat_back_id")
		}
		return int64(content.Id), err
	}
	id, err = orm.NewOrm().Insert(data)
	return
}

//直接保存内容，根据ID判断是否更新
func (c *Content) SaveUpdate(data *Content) (id int64, err error) {
	if data.Id > 0 {
		_, err = c.UpdateContent(data, "title", "content", "source_url", "sourceurl_crc32", "author_id", "genre_id", "creat_back_id", "sort")
		id = int64(data.Id)
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}

//根据标题搜索诗词信息
func (c *Content) GetByTitleAuthorId(title string, authorId int64) (data Content, err error) {
	_, err = orm.NewOrm().QueryTable(TableContent).Filter("title", title).Filter("author_id", authorId).All(&data, "id", "content")
	return
}

//根据URL的crc32值查询
func (c *Content) GetContentByCrc32(crc32 uint32) (data Content, err error) {
	_, err = orm.NewOrm().QueryTable(TableContent).Filter("sourceurl_crc32", crc32).All(&data, "id")
	return
}

//更新数据
func (c *Content) UpdateContent(data *Content, col ...string) (id int64, err error) {
	id, err = orm.NewOrm().Update(data, col...)
	return
}
