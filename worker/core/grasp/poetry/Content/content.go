package Content

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Helper"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"regexp"
	"strconv"
	"strings"
)

//诗文详情模块 抓取诗文详情
//https://so.gushiwen.org/shiwenv_73add8822103.aspx
type Content struct {
}

func NewContent() *Content {
	return &Content{}
}

//传过来一个诗词详情页的URL(/shiwenv_73add8822103.aspx)，获取数据并保存诗词详情数据
func (c *Content) GraspContentSaveData(detailUrl string, params interface{}) {
	var (
		bytes  []byte
		err    error
		urlKey uint32
	)
	url := detailUrl
	if strings.Contains(detailUrl, "http") == false {
		url = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(detailUrl, "/")
	}
	urlKey = tools.Crc32(url)
	if ok := tools.NewLock().ExistsKey(urlKey); ok {
		logrus.Infoln("GraspContentSaveData:", url, "重复请求....")
		return
	}
	tools.NewLock().AddKey(urlKey)
	if bytes, err = Helper.GetContentHtml(url); err != nil {
		data.G_GraspResult.PushError(err)
		tools.NewLock().DelKey(urlKey)
		return
	}
	content := c.FindDocumentData(bytes)
	content.SourceUrl = url

	/*
		sendData := &define.ParseData{
			Data:      &content,
			Params:    params,
			ParseFunc: data.NewContentStore().LoadPoetryContentData,
		}
		data.G_GraspResult.SendParseData(sendData)
	*/
	logrus.Infof("content:%+v\n", content)
	data.NewContentStore().LoadPoetryContentData(&content, params)
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
	poetryContent.Author = Helper.GetAuthorData(query)
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
		titleMap    map[string]interface{}
	)
	titleMap = make(map[string]interface{})
	query.Find(".left>.sons").Each(func(i int, selection *goquery.Selection) {
		var buf bytes.Buffer
		idStr, exists := selection.Attr("id")
		attr, ok = selection.Find("a").Attr("href")
		title := selection.Find(".contyishang>div>h2").Text()
		_, ok := titleMap[title]
		if exists == true && !ok && len(title) > 0 {
			titleMap[title] = 1
			selection.Find(".contyishang>p").Each(func(i int, selection *goquery.Selection) {
				if html, e := selection.Html(); e == nil {
					buf.WriteString("<p>" + html + "</p>")
				}
			})
			//翻译
			if strings.Contains(idStr, "fanyi") && !strings.Contains(idStr, "fanyiquan") {
				id = strings.TrimLeft(idStr, "fanyi")
				notesUrl = config.G_Conf.GuShiWenPoetryUrl + "nocdn/ajaxfanyi.aspx?id=" + id
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
				apprecUrl = config.G_Conf.GuShiWenPoetryUrl + "nocdn/ajaxshangxi.aspx?id=" + appRecId
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
		if exists == false && !ok && len(title) > 0 {
			content := &define.PoetryContentData{
				TransId:    0,
				Content:    "",
				Introd:     "",
				HtmlSrcUrl: "",
				Title:      title,
				Sort:       i,
			}
			if len(attr) > 0 && strings.Contains(attr, "javascript:PlayFanyi") {
				idStr := strings.Replace(attr, "javascript:PlayFanyi(", "", -1)
				idStr = strings.TrimRight(idStr, ")")
				trId, _ := strconv.Atoi(idStr)
				content.TransId = trId
				content.PlaySrcUrl = config.G_Conf.GuShiWenPoetryUrl + "fanyiplay.aspx?id=" + idStr
				content.PlayUrl = config.G_Conf.GushiwenSongUrl + "machine/fanyi/" + idStr + "/ok.mp3"
			}
			if len(attr) > 0 && strings.Contains(attr, "javascript:PlayShangxi") {
				idStr := strings.Replace(attr, "javascript:PlayShangxi(", "", -1)
				idStr = strings.TrimRight(idStr, ")")
				appId, _ := strconv.Atoi(idStr)
				content.AppRecId = appId
				content.PlayUrl = config.G_Conf.GushiwenSongUrl + "machine/shangxi/" + idStr + "/ok.mp3"
				content.PlaySrcUrl = config.G_Conf.GuShiWenPoetryUrl + "/shangxiplay.aspx?id=" + idStr
			}
			content.Content, _ = selection.Find(".contyishang").Html()
			if len(content.Content) > 0 {
				content.Content = tools.TrimDivHtml(content.Content)
				mustCompile := regexp.MustCompile(`(?msU)<div.*>.*</div>`)
				content.Content = mustCompile.ReplaceAllString(content.Content, "")
			}
			notesData = append(notesData, content)
		}

	})
	//https://so.gushiwen.org/shiwenv_58313be2d918.aspx
	if isTransData == false && !ok {
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

//通过诗文详情数据 获取作者信息 https://so.gushiwen.org/shiwenv_4c5705b99143.aspx
func (c *Content) GetAuthorContentData(poetry *define.PoetryAuthorList) (author *define.PoetryAuthorDetail) {
	var (
		bytes  []byte
		query  *goquery.Document
		err    error
		urlKey uint32
	)
	url := poetry.PoetrySourceUrl
	if strings.Contains(url, "http") == false {
		url = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(poetry.PoetrySourceUrl, "/")
	}
	//过虑重复请求，有可能多个分类下有同一个作者， 这里只取一次作者信息
	urlKey = tools.Crc32(url)
	if ok := tools.NewLock().ExistsKey(urlKey); ok {
		logrus.Infoln(url, "重复请求....")
		return
	}
	tools.NewLock().AddKey(urlKey)
	if bytes, err = Helper.GetContentHtml(url); err != nil {
		data.G_GraspResult.PushError(err)
		tools.NewLock().DelKey(urlKey)
		return
	}
	if query, err = tools.NewDocumentFromReader(string(bytes)); err != nil {
		data.G_GraspResult.PushError(err)
		tools.NewLock().DelKey(urlKey)
		return
	}
	author = Helper.GetAuthorData(query)
	return author
}
