package Content

import (
	"bytes"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Author"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"strconv"
	"strings"
	"sync"
)

//诗文详情模块 抓取诗文详情
//https://so.gushiwen.org/shiwenv_73add8822103.aspx
type Content struct {
	contentMap sync.Map
	saveMap    sync.Map
}

func NewContent() *Content {
	return &Content{}
}

//通过诗文分类 抓取诗文详情数据 https://so.gushiwen.org/shiwenv_4c5705b99143.aspx
func (c *Content) GraspContentData(poetry *define.PoetryAuthorList) {
	var (
		bytes          []byte
		query          *goquery.Document
		err            error
		urlKey         uint32
		author         *define.PoetryAuthorDetail
		authorChan     chan bool
		authorListChan chan bool
	)
	url := config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(poetry.PoetrySourceUrl, "/")
	//过虑重复请求，有可能多个分类下有同一个作者， 这里只取一次作者信息
	urlKey = tools.Crc32(poetry.PoetrySourceUrl)
	if ok := NewLock().ExistsKey(urlKey); ok {
		logrus.Infoln(url, "重复请求....")
		return
	}
	NewLock().AddKey(urlKey)
	if bytes, err = c.GetContentSource(url); err != nil {
		data.G_GraspResult.PushError(err)
		NewLock().DelKey(urlKey)
		return
	}
	if query, err = tools.NewDocumentFromReader(string(bytes)); err != nil {
		data.G_GraspResult.PushError(err)
		NewLock().DelKey(urlKey)
		return
	}
	author = c.getAuthorData(query)
	authorChan, authorListChan = make(chan bool), make(chan bool)
	//发送获取作者详情信息请求
	go Author.NewAuthor().SetAuthorAttr(author).GraspAuthorDetail(author.AuthorSrcUrl, authorChan)
	//发送获取作者诗词列表所有数据请求
	go Author.NewAuthor().SetAuthorAttr(author).GraspAuthorPoetryList(author.AuthorContentUrl, authorListChan)
	<-authorChan
	<-authorListChan
	return
}

//传过来一个诗词详情页的URL(/shiwenv_73add8822103.aspx)，获取数据并保存诗词详情数据
func (c *Content) GraspContentSaveData(detailUrl string, params []interface{}) {
	var (
		bytes  []byte
		err    error
		urlKey uint32
	)
	url := config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(detailUrl, "/")
	urlKey = tools.Crc32(url)
	if ok := NewLock().ExistsKey(urlKey); ok {
		logrus.Infoln("GraspContentSaveData:", url, "重复请求....")
		return
	}
	if bytes, err = c.GetContentSource(url); err != nil {
		data.G_GraspResult.PushError(err)
		NewLock().DelKey(urlKey)
		return
	}
	content := c.FindDocumentData(bytes)
	sendData := &define.ParseData{
		Data:      &content,
		Params:    params,
		ParseFunc: data.NewContentStore().LoadPoetryContentData,
	}
	data.G_GraspResult.SendParseData(sendData)
	NewLock().AddKey(urlKey)
}

//获取诗文详情数据
func (c *Content) GetContentSource(url string) (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		bytes, err = c.getTestFile()
	} else {
		bytes, err = base.GetHtml(url)
	}
	return
}

//goquery 查找数据
func (c *Content) FindDocumentData(html []byte) (poetryContent define.PoetryContent) {
	var (
		query *goquery.Document
		err   error
	)
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	src := ".left>.sons>.cont"
	poetryContent.Title = query.Find(src + ">h1").Text()
	poetryContent.Content = query.Find(src + ">.contson").Eq(0).Text()
	poetryContent.CategoryList = c.getCategory(query)
	poetryContent.Author = c.getAuthorData(query)
	poetryContent.Detail = c.getNotesData(query)
	poetryContent.CreativeBackground = c.getCreativeBack(query) //创作背景
	return poetryContent
}

//获取译文及注释与赏析数据  [这里还要抓内容简介，暂时没做... 抓诗词的时候加上]
func (c *Content) getNotesData(query *goquery.Document) (notesData []*define.PoetryContentData) {
	var (
		notesUrl    string
		apprecUrl   string
		appRecId    string
		id          string
		ok          bool
		attr        string
		isTransData bool
		htmlStr     string
	)
	query.Find(".left>.sons").Each(func(i int, selection *goquery.Selection) {
		var buf bytes.Buffer
		idStr, exists := selection.Attr("id")
		attr, ok = selection.Find("a").Attr("href")
		title := selection.Find(".contyishang>div>h2").Text()
		if exists == true {
			selection.Find(".contyishang>p").Each(func(i int, selection *goquery.Selection) {
				if html, e := selection.Html(); e == nil {
					buf.WriteString("<p>" + html + "</p>")
				}
			})
			//翻译
			if strings.Contains(idStr, "fanyi") && !strings.Contains(idStr, "fanyiquan") {
				id = strings.TrimLeft(idStr, "fanyi")
				notesUrl = config.G_Conf.GuShiWenPoetryUrl + "shiwen2017/ajaxfanyi.aspx?id=" + id
				if bytes, err := base.GetHtml(notesUrl); err == nil {
					trId, _ := strconv.Atoi(id)
					content := &define.PoetryContentData{
						TransId:    trId,
						Content:    string(bytes),
						Introd:     buf.String(),
						HtmlSrcUrl: notesUrl,
						Title:      title,
						Sort:       i,
					}
					if len(attr) > 0 && strings.Contains(attr, "javascript:PlayFanyi") {
						content.PlaySrcUrl = config.G_Conf.GuShiWenPoetryUrl + "fanyiplay.aspx?id=" + id
						content.PlayUrl = config.G_Conf.GushiwenSongUrl + "machine/fanyi/" + id + "/ok.mp3"
					}
					buf.Reset()
					notesData = append(notesData, content)
					isTransData = true
				}
			}
			//赏析
			if strings.Contains(idStr, "shangxi") && !strings.Contains(idStr, "shangxiquan") {
				appRecId = strings.TrimLeft(idStr, "shangxi")
				apprecUrl = config.G_Conf.GuShiWenPoetryUrl + "shiwen2017/ajaxshangxi.aspx?id=" + appRecId
				if bytes, err := base.GetHtml(apprecUrl); err == nil {
					appId, _ := strconv.Atoi(appRecId)
					content := &define.PoetryContentData{
						AppRecId:   appId,
						Content:    string(bytes),
						Introd:     buf.String(),
						HtmlSrcUrl: apprecUrl,
						Title:      title,
						Sort:       i,
					}
					if len(attr) > 0 && strings.Contains(attr, "javascript:PlayShangxi") {
						content.PlayUrl = config.G_Conf.GushiwenSongUrl + "machine/shangxi/" + appRecId + "/ok.mp3"
						content.PlaySrcUrl = config.G_Conf.GuShiWenPoetryUrl + "/shangxiplay.aspx?id=" + appRecId
					}
					buf.Reset()
					notesData = append(notesData, content)
				}
			}
		}
	})
	//https://so.gushiwen.org/shiwenv_58313be2d918.aspx
	if isTransData == false {
		conty1 := query.Find(".left>.sons>.contyishang").Eq(0)
		conty1.Find("p").Each(func(i int, selection *goquery.Selection) {
			if html, err := selection.Html(); err == nil {
				htmlStr += html
			}
		})
		if attr, ok = conty1.Find("a").Attr("href"); ok {
			attr = strings.TrimLeft(attr, "javascript:PlayFanyi(")
			attr = strings.TrimRight(attr, ")")
			trId, _ := strconv.Atoi(attr)
			title := conty1.Find("div>h2").Text()
			content := &define.PoetryContentData{
				TransId:    trId,
				Content:    htmlStr,
				Sort:       1,
				Title:      title,
				PlaySrcUrl: config.G_Conf.GuShiWenPoetryUrl + "fanyiplay.aspx?id=" + attr,
				PlayUrl:    config.G_Conf.GushiwenSongUrl + "machine/fanyi/" + attr + "/ok.mp3",
			}
			notesData = append(notesData, content)
		}
	}
	return
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

//获取作者头像和总诗词数
func (c *Content) getAuthorData(query *goquery.Document) (author *define.PoetryAuthorDetail) {
	authorImg, _ := query.Find(".sonspic>.cont>.divimg>a>img").Attr("src")
	text := query.Find(".sonspic>.cont>p").Eq(1).Find("a").Text()
	authorSrcUrl, _ := query.Find(".sonspic>.cont>p").Eq(1).Find("a").Attr("href")
	text = strings.TrimRight(text, "篇诗文")
	text = strings.TrimLeft(text, " ► ")
	total, _ := strconv.Atoi(text)
	author = &define.PoetryAuthorDetail{
		AuthorImgUrl:      authorImg,
		AuthorContentUrl:  authorSrcUrl,
		AuthorTotalPoetry: total,
	}
	src := ".left>.sons>.cont"
	author.DynastyName = query.Find(src + ">.source>a").Eq(0).Text()
	author.DynastyUrl, _ = query.Find(src + ">.source>a").Eq(0).Attr("href")
	author.AuthorName = query.Find(src + ">.source>a").Eq(1).Text()
	author.AuthorSrcUrl, _ = query.Find(src + ">.source>a").Eq(1).Attr("href")
	return
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
func (c *Content) getTestFile() (bytes []byte, err error) {
	dir, _ := os.Getwd()
	file := dir + "/content1.html"
	if ret, _ := tools.PathExists(file); ret == true {
		return tools.ReadFile(file)
	}
	return nil, errors.New(file + "file is not exists")
}
