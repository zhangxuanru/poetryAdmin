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
	conty1 := query.Find(".left>.sons>.contyishang").Eq(0)
	conty1.Find("p").Each(func(i int, selection *goquery.Selection) {
		text := selection.Text()
		if text != "注释" {
			poetryContent.Translation += text //译文
		} else {
			poetryContent.Notes += text //注释
		}
	})
	//poetryContent.Translation, _ = conty1.Find("p").Eq(0).Html() //译文
	//poetryContent.Notes, _ = conty1.Find("p").Eq(1).Html()       //注释
	//赏析
	if poetryContent.Appreciation, err = c.getAppreciation(query); err != nil {
		poetryContent.Appreciation = ""
		data.G_GraspResult.PushError(err)
	}
	poetryContent.CreativeBackground = query.Find(".left>.sons").Eq(4).Find(".contyishang>p").Eq(0).Text()

	logrus.Infof("%+v", poetryContent)

}

//获取赏析数据
func (c *Content) getAppreciation(query *goquery.Document) (body string, err error) {
	var (
		appRecId        string
		ok              bool
		appreciationUrl string
	)
	appRecId, ok = query.Find(".sons").Eq(2).Attr("id")
	appRecId = strings.TrimLeft(appRecId, "shangxi")
	if ok == false || appRecId == "" {
		body = query.Find(".sons").Eq(2).Find(".contyishang>p").Eq(0).Text()
		body += query.Find(".sons").Eq(2).Find(".contyishang>p").Eq(1).Text()
		return body, nil
	}
	return
	appreciationUrl = config.G_Conf.GuShiWenPoetryUrl + "shiwen2017/ajaxshangxi.aspx?id=" + appRecId
	bytes, err := base.GetHtml(appreciationUrl)
	body = "html:" + string(bytes)
	return body, err
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
