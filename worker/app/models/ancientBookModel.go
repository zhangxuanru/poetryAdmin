/*
@Time : 2019/9/3 18:52
@Author : zxr
@File : ancientBookModel
@Software: GoLand
*/
package models

import "github.com/astaxie/beego/orm"

var TableAncientBook = "poetry_ancient_book"

//poetry_ancient_book 古籍-书名表
type AncientBook struct {
	Id               int64  `orm:"column(id);auto"`
	CatId            int    `orm:"column(cat_id)"`
	BookTitle        string `orm:"column(book_title)"`
	BookIntroduction string `orm:"column(book_introduction)"`
	LinkUrl          string `orm:"column(link_url)"`
	LinkUrlCrc32     uint32 `orm:"column(link_url_crc32)"`
	SongUrl          string `orm:"column(song_url)"`
	SongFilePath     string `orm:"column(song_file_path)"`
	SongSrcUrl       string `orm:"column(song_src_url)"`
	FamousTotal      int    `orm:"column(famous_total)"`
	CoverChart       string `orm:"column(cover_chart)"`
	CoverChartPath   string `orm:"column(cover_chart_path)"`
	AddDate          int64  `orm:"column(add_date)"`
}

func init() {
	orm.RegisterModel(new(AncientBook))
}

func (a *AncientBook) TableName() string {
	return TableAncientBook
}

func NewAncientBook() *AncientBook {
	return new(AncientBook)
}

//保存书籍信息
func (a *AncientBook) SaveBook(data *AncientBook) (id int64, err error) {
	var book AncientBook
	if book, err = a.GetBookByTitleAndCatId(data.BookTitle, data.CatId); err != nil {
		return 0, err
	}
	if book.Id > 0 {
		//	_, err = a.UpdateAllData(data)
		return int64(book.Id), nil
	} else {
		id, err = orm.NewOrm().Insert(data)
	}
	return
}

//根据标题和分类查询书籍数据
func (a *AncientBook) GetBookByTitleAndCatId(title string, catId int) (data AncientBook, err error) {
	_, err = orm.NewOrm().QueryTable(TableAncientBook).Filter("cat_id", catId).Filter("book_title", title).All(&data, "id")
	return
}

//根据标题和urlcrc32值查询
func (a *AncientBook) GetBookByTitleANDUrlCrc32(title string, urlCrc uint32) (data AncientBook, err error) {
	_, err = orm.NewOrm().QueryTable(TableAncientBook).Filter("title", title).Filter("link_url_crc32", urlCrc).All(&data, "id")
	return
}

//更新所有数据
func (a *AncientBook) UpdateAllData(data *AncientBook) (id int64, err error) {
	id, err = orm.NewOrm().Update(data)
	return
}

//更新声音保存的本地路径
func (a *AncientBook) UpdateSongPath(bookId int64, songFilePath string) (id int64, err error) {
	data := &AncientBook{
		Id:           bookId,
		SongFilePath: songFilePath,
	}
	id, err = orm.NewOrm().Update(data, "song_file_path")
	return
}

//更新封面图地址本地路径
func (a *AncientBook) UpdateCoverPath(bookId int64, coverChartPath string) (id int64, err error) {
	data := &AncientBook{
		Id:             bookId,
		CoverChartPath: coverChartPath,
	}
	id, err = orm.NewOrm().Update(data, "cover_chart_path")
	return
}
