/*
@Time : 2019/9/2 14:37
@Author : zxr
@File : ancient_test
@Software: GoLand
*/
package test

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/ancient/Action"
	"poetryAdmin/worker/core/grasp/ancient/Entrance"
	"poetryAdmin/worker/core/grasp/ancient/Parser"
	"testing"
	"time"
)

//测试古籍
func TestAncient(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	Entrance.NewGrab().Exec()
	time.Sleep(60 * time.Second)
}

//测试抓取目录
func TestAncientCatLog(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	bookCoverList := []*define.GuWenBookCover{
		&define.GuWenBookCover{
			Title:   "周礼",
			LinkUrl: "https://so.gushiwen.org/guwen/book_24.aspx",
		},
	}
	category := &define.GuWenCategory{
		CategoryName: "词曲类",
		LinkUrl:      "/guwen/Default.aspx?p=1&type=%e8%af%8d%e6%9b%b2%e7%b1%bb",
	}
	Action.NewCataLog().LoadBookCoverList(bookCoverList, category)

	time.Sleep(20 * time.Second)
}

//抓取古文内容
func TestAncientContent(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	book := &define.BookCatalogue{
		BookTitle:   "周礼",
		BookLinkUrl: "https://so.gushiwen.org/guwen/book_24.aspx",
		CatalogList: []define.CataLog{
			{
				CatalogTitle: "大司马",
				LinkUrl:      "/guwen/bookv_3218.aspx",
			},
			{
				CatalogTitle: "大宗伯",
				LinkUrl:      "/guwen/bookv_3208.aspx",
			},
			{
				CatalogTitle: "大宰",
				LinkUrl:      "/guwen/bookv_3187.aspx",
			},
		},
	}
	Action.NewContent().LoadBookCatalogue(book)

	time.Sleep(5 * time.Second)
}

func TestA(t *testing.T) {
	html, _ := Action.NewCataLog().GetCataLogHtml("https://so.gushiwen.org/guwen/book_20.aspx")
	catalogData, _ := Parser.ParseGuWenCataLog(html)
	for _, v := range catalogData {
		logrus.Infof("%+v\n", v)
	}
}
