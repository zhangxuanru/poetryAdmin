package models

import "github.com/astaxie/beego/orm"

var TableContent = "poetry_content"

//poetry_content 诗词表
type Content struct {
	Id         int    `orm:"column(id);auto"`
	Title      string `orm:"column(title)"`
	Content    string `orm:"column(content)"`
	AuthorId   int64  `orm:"column(author_id)"`
	SourceUrl  string `orm:"column(source_url)"`
	AddDate    int64  `orm:"column(add_date)"`
	UpdateDate int64  `orm:"column(update_date)"`
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
	if content, err = c.GetByTitleAuthorId(data.Title, data.AuthorId); err != nil {
		return 0, err
	}
	if content.Id > 0 {
		data.Id = content.Id
		_, _ = c.UpdateContent(data, "update_date")
		return int64(content.Id), nil
	}
	id, err = orm.NewOrm().Insert(data)
	return
}

//根据标题搜索诗词信息
func (c *Content) GetByTitleAuthorId(title string, authorId int64) (data Content, err error) {
	_, err = orm.NewOrm().QueryTable(TableContent).Filter("title", title).Filter("author_id", authorId).All(&data, "id", "content", "update_date")
	return
}

//更新数据
func (c *Content) UpdateContent(data *Content, col ...string) (id int64, err error) {
	id, err = orm.NewOrm().Update(data, col...)
	return
}
