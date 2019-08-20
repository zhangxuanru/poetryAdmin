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
//https://so.gushiwen.org/shiwenv_73add8822103.aspx
type Content struct {
}

func NewContent() *Content {
	return &Content{}
}

//通过诗文分类 抓取诗文详情数据
func (c *Content) GraspCategoryData(poetry *define.PoetryAuthorList) {
	url := config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(poetry.PoetrySourceUrl, "/")
	bytes, err := c.GetContentSource(url)
	if err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	c.FindDocumentData(bytes)
	return
}

//获取诗文详情数据
func (c *Content) GetContentSource(url string) (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		bytes, err = c.CateTestFile()
	} else {
		bytes, err = base.GetHtml(url)
	}
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
	poetryContent.CategoryList = c.getCategory(query)
	notes, _ := c.getNotes(query) //译文及注释
	poetryContent.Notes = strings.Join(notes, "#")
	appreciation, _ := c.getAppreciation(query)
	poetryContent.Appreciation = strings.Join(appreciation, "#") //赏析
	poetryContent.CreativeBackground = c.getCreativeBack(query)  //创作背景
	logrus.Infof("%+v", poetryContent)
}

//获取译文  译文及注释放在一起
func (c *Content) getNotes(query *goquery.Document) (body []string, err error) {
	var (
		notesUrl string
		id       string
		htmlStr  string
	)
	//-----
	query.Find(".left>.sons").Each(func(i int, selection *goquery.Selection) {
		idStr, exists := selection.Attr("id")
		if exists == true && !strings.Contains(idStr, "fanyiquan") && strings.Contains(idStr, "fanyi") {
			id = strings.TrimLeft(idStr, "fanyi")
			notesUrl = config.G_Conf.GuShiWenPoetryUrl + "shiwen2017/ajaxfanyi.aspx?id=" + id
			if bytes, err := base.GetHtml(notesUrl); err == nil {
				data.G_GraspResult.PushError(err)
				body = append(body, string(bytes))
			}
		}
	})
	if len(body) == 0 {
		conty1 := query.Find(".left>.sons>.contyishang").Eq(0)
		conty1.Find("p").Each(func(i int, selection *goquery.Selection) {
			if html, err := selection.Html(); err == nil {
				htmlStr += html
			}
		})
		body = append(body, htmlStr)
	}
	return
	//-----

	//id, ok = query.Find(".left>.sons").Eq(1).Attr("id")
	//id = strings.TrimLeft(id, "fanyi")
	//if ok == false || id == "" {
	//	conty1 := query.Find(".left>.sons>.contyishang").Eq(0)
	//	conty1.Find("p").Each(func(i int, selection *goquery.Selection) {
	//		if html, err := selection.Html(); err == nil {
	//			body += html
	//		}
	//	})
	//	return
	//}
	//notesUrl = config.G_Conf.GuShiWenPoetryUrl + "shiwen2017/ajaxfanyi.aspx?id=" + id
	//bytes, err = base.GetHtml(notesUrl)
	//if err != nil {
	//	return "", err
	//}
	//body = string(bytes)
	//return
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

//获取诗的分类数据
func (c *Content) getCategory(query *goquery.Document) (categoryList []*define.TextHrefFormat) {
	query.Find(".left>.sons>.tag").Eq(0).Find("a").Each(func(i int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		format := &define.TextHrefFormat{
			Text:         selection.Text(),
			Href:         href,
			ShowPosition: 1,
		}
		categoryList = append(categoryList, format)
	})
	return categoryList
}

//获取诗的创作背景
func (c *Content) getCreativeBack(query *goquery.Document) (body string) {
	query.Find(".left>.sons>.contyishang").Each(func(i int, selection *goquery.Selection) {
		text := selection.Find("div>h2").Text()
		if text == "创作背景" {
			body, _ = selection.Find("p").Html()
		}
	})
	return
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
