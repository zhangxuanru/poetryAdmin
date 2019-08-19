package Content

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"strings"
)

//诗文详情模块 抓取诗文详情
type Content struct {
}

func NewContent() *Content {
	return &Content{}
}

//通过诗文分类 抓取诗文详情数据
func (c *Content) GraspCategoryData(poetry *define.PoetryAuthorList) {
	url := config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(poetry.PoetrySourceUrl, "/")
	c.GetContentSource(url)
}

//获取诗文详情数据
func (c *Content) GetContentSource(url string) {
	var (
		bytes []byte
		err   error
	)
	if config.G_Conf.Env == define.TestEnv {
		bytes, err = c.CateTestFile()
	} else {
		bytes, err = base.GetHtml(url)
	}
	if err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	c.FindDocumentData(bytes)
	return
}

//goquery 查找数据
func (c *Content) FindDocumentData(html []byte) {
	var (
		query         *goquery.Document
		poetryContent define.PoetryContent
		err           error
	)
	query, err = tools.NewDocumentFromReader(string(html))
	if err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	src := ".left>.sons>.cont"
	poetryContent.Title = query.Find(src + ">h1").Text()
	poetryContent.DynastyName = query.Find(src + ">.source>a").Eq(0).Text()
	poetryContent.DynastyUrl, _ = query.Find(src + ">.source>a").Eq(0).Attr("href")
	poetryContent.AuthorName = query.Find(src + ">.source>a").Eq(1).Text()
	poetryContent.AuthorSrcUrl, _ = query.Find(src + ">.source>a").Eq(1).Attr("href")
	poetryContent.Content = query.Find(src + ">.contson").Eq(0).Text()
	query.Find(".left>.sons>.tag").Eq(0).Find("a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		format := &define.TextHrefFormat{
			Text:         selection.Text(),
			Href:         href,
			ShowPosition: 1,
		}
		poetryContent.CategoryList = append(poetryContent.CategoryList, format)
	})
	poetryContent.Notes, _ = c.getNotes(query) //译文及注释
	appreciation, _ := c.getAppreciation(query)
	poetryContent.Appreciation = strings.Join(appreciation, "#") //赏析

	//创作背景明天继续.... 这里没想好怎么处理
	poetryContent.CreativeBackground = query.Find(".left>.sons").Eq(4).Find(".contyishang>p").Eq(0).Text()

	logrus.Infof("%+v", poetryContent)

}

//获取译文  译文及注释放在一起
func (c *Content) getNotes(query *goquery.Document) (body string, err error) {
	var (
		notesUrl string
		id       string
		bytes    []byte
		ok       bool
	)
	id, ok = query.Find(".left>.sons").Eq(1).Attr("id")
	id = strings.TrimLeft(id, "fanyi")
	if ok == false || id == "" {
		conty1 := query.Find(".left>.sons>.contyishang").Eq(0)
		conty1.Find("p").Each(func(i int, selection *goquery.Selection) {
			if html, err := selection.Html(); err == nil {
				body += html
			}
		})
		return
	}
	notesUrl = config.G_Conf.GuShiWenPoetryUrl + "shiwen2017/ajaxfanyi.aspx?id=" + id
	bytes, err = base.GetHtml(notesUrl)
	if err != nil {
		return "", err
	}
	body = string(bytes)
	return
}

//获取赏析数据
func (c *Content) getAppreciation(query *goquery.Document) (body []string, err error) {
	var (
		appRecId        string
		appreciationUrl string
	)
	query.Find(".sons").Each(func(i int, selection *goquery.Selection) {
		idStr, exists := selection.Attr("id")
		if exists == true && !strings.Contains(idStr, "shangxiquan") && strings.Contains(idStr, "shangxi") {
			appRecId = strings.TrimLeft(idStr, "shangxi")
			appreciationUrl = config.G_Conf.GuShiWenPoetryUrl + "shiwen2017/ajaxshangxi.aspx?id=" + appRecId
			bytes, _ := base.GetHtml(appreciationUrl)
			body = append(body, string(bytes))
		}
	})
	return body, nil
}

//读取测试文件内容
func (c *Content) CateTestFile() (bytes []byte, err error) {
	dir, _ := os.Getwd()
	file := dir + "/content1.html"
	if ret, _ := tools.PathExists(file); ret == true {
		return tools.ReadFile(file)
	}
	return nil, errors.New(file + "file is not exists")
}
