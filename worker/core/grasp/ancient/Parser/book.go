/*
@Time : 2019/9/2 18:22
@Author : zxr
@File : book
@Software: GoLand
*/
package Parser

import (
	"github.com/PuerkitoBio/goquery"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"regexp"
	"strconv"
	"strings"
)

//解析分类下书籍信息
func ParseGuWenCategoryBook(html []byte) (bookCoverList []*define.GuWenBookCover, err error) {
	var (
		query *goquery.Document
	)
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	query.Find(".sonspic").Each(func(i int, selection *goquery.Selection) {
		var bookCover define.GuWenBookCover
		bookCover.CoverChart, _ = selection.Find(".cont>.divimg>a>img").Attr("src")
		bookCover.Title = selection.Find("p").Eq(0).Find("a>b").Text()
		bookCover.LinkUrl, _ = selection.Find("p").Eq(0).Find("a").Eq(0).Attr("href")
		songUrl, _ := selection.Find("p").Eq(0).Find("a").Eq(1).Attr("href")
		if len(songUrl) > 0 {
			songUrl = strings.TrimLeft(songUrl, "javascript:PlayBook(")
			songId := strings.TrimRight(songUrl, ")")
			bookCover.SongUrl = "https://song.gushiwen.org/machine/book/" + songId + "/ok.mp3"
			bookCover.SongSrcUrl = "https://so.gushiwen.org/guwen/bookplay.aspx?id=" + songId
		}
		bookCover.Introduction, _ = selection.Find("p").Eq(1).Html()
		famousText := selection.Find("p").Eq(1).Find("a").Text()
		if len(famousText) > 0 {
			famousText = strings.TrimLeft(famousText, "► ")
			famousText = strings.TrimRight(famousText, "条名句")
			bookCover.FamousTotal, err = strconv.Atoi(famousText)
			compile := regexp.MustCompile(`<a(.*)</a>`)
			bookCover.Introduction = compile.ReplaceAllString(bookCover.Introduction, "")
		}
		bookCover.Title = strings.TrimSpace(bookCover.Title)
		bookCover.Introduction = strings.TrimSpace(bookCover.Introduction)
		bookCoverList = append(bookCoverList, &bookCover)
	})
	return
}
