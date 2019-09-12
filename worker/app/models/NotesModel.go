package models

import "github.com/astaxie/beego/orm"

var TableNotes = "poetry_detail_notes"

//poetry_detail_notes 诗词详情内容表
type Notes struct {
	Id         int    `orm:"column(id);auto"`
	Title      string `orm:"column(title)"`
	Content    string `orm:"column(content)"`
	PlayUrl    string `orm:"column(play_url)"`
	PlaySrcUrl string `orm:"column(play_src_url)"`
	HtmlSrcUrl string `orm:"column(html_src_url)"`
	Type       int    `orm:"column(type)"`
	Introd     string `orm:"column(introd)"`
	FileName   string `orm:"column(file_name)"`
	AddDate    int64  `orm:"column(add_date)"`
	UpdateDate int64  `orm:"column(update_date)"`
}

func init() {
	orm.RegisterModel(new(Notes))
}

func NewNotes() *Notes {
	return new(Notes)
}

func (n *Notes) TableName() string {
	return TableNotes
}

//保存内容
func (n *Notes) SaveNotes(data *Notes) (id int64, err error) {
	if data.Id > 0 {
		_, err = n.UpdateNotes(data)
		id = int64(data.Id)
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}

//更新数据
func (n *Notes) UpdateNotes(data *Notes, col ...string) (id int64, err error) {
	if len(col) == 0 {
		col = []string{"play_url", "play_src_url", "html_src_url", "type", "title", "introd", "content", "update_date"}
	}
	id, err = orm.NewOrm().Update(data, col...)
	return
}
