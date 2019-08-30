/*
@Time : 2019/8/27 18:02
@Author : zxr
@File : Content
@Software: GoLand
*/
package Helper

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"strconv"
	"strings"
)

//获取作者头像和总诗词数
func GetAuthorData(query *goquery.Document) (author *define.PoetryAuthorDetail) {
	authorImg, _ := query.Find(".sonspic>.cont>.divimg>a>img").Attr("src")
	if len(authorImg) == 0 {
		authorImg, _ = query.Find(".sonspic>.cont>.divimg>img").Attr("src")
	}
	text := query.Find(".sonspic>.cont>p").Eq(1).Find("a").Text()
	if len(text) == 0 {
		text = query.Find(".sonspic>.cont>p").Eq(0).Find("a").Text()
	}
	authorSrcUrl, _ := query.Find(".sonspic>.cont>p").Eq(1).Find("a").Attr("href")
	if len(authorSrcUrl) == 0 {
		authorSrcUrl, _ = query.Find(".sonspic>.cont>p").Eq(0).Find("a").Attr("href")
	}
	text = strings.TrimRight(text, "篇诗文")
	text = strings.TrimLeft(text, " ► ")
	total, _ := strconv.Atoi(text)
	author = &define.PoetryAuthorDetail{
		AuthorImgUrl:      authorImg,
		AuthorContentUrl:  strings.TrimSpace(authorSrcUrl),
		AuthorTotalPoetry: total,
	}
	src := ".left>.sons>.cont"
	author.DynastyName = query.Find(src + ">.source>a").Eq(0).Text()
	author.DynastyUrl, _ = query.Find(src + ">.source>a").Eq(0).Attr("href")
	author.AuthorName = query.Find(src + ">.source>a").Eq(1).Text()
	author.AuthorSrcUrl, _ = query.Find(src + ">.source>a").Eq(1).Attr("href")
	return
}

//获取诗文详情html数据
func GetContentHtml(url string) (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		bytes, err = GetTestContentFile()
	} else {
		bytes, err = base.GetHtml(url)
	}
	return
}

//读取测试文件内容
func GetTestContentFile() (bytes []byte, err error) {
	dir, _ := os.Getwd()
	file := dir + "/content3.html"
	if ret, _ := tools.PathExists(file); ret == true {
		return tools.ReadFile(file)
	}
	return nil, errors.New(file + "file is not exists")
}
