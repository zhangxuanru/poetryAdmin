package Author

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Content"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"qiniupkg.com/x/errors.v7"
	"strconv"
	"strings"
	"time"
)

type PoetryLinkMap map[int]define.LinkStr //诗词标题与链接信息

//作者模块  抓取作者
type Author struct {
	SourceAuthor *define.PoetryAuthorDetail
	Html         []byte
	query        *goquery.Document
	Url          string
}

func NewAuthor() *Author {
	return &Author{}
}

//发送获取作者详情和 诗词列表的请求
func (a *Author) SendGraspAuthorDataReq(author *define.PoetryAuthorDetail, srcUrl string) {
	var (
		authorChan     chan bool
		authorListChan chan bool
	)

	logrus.Infoln("AuthorContentUrl:", author.AuthorContentUrl)
	logrus.Infof("author:%+v\n------------\n", author)

	authorChan, authorListChan = make(chan bool), make(chan bool)
	//发送获取作者详情信息请求
	a.SetAuthorAttr(author).GraspAuthorDetail(author.AuthorSrcUrl, authorChan)
	<-authorChan
	//发送获取作者诗词列表所有数据请求
	if len(author.AuthorContentUrl) > 0 {
		a.SetAuthorAttr(author).GraspAuthorPoetryList(author.AuthorContentUrl, authorListChan)
		<-authorListChan
	} else {
		logrus.Infoln("srcUrl:", srcUrl)
		logrus.Infof("%s---%s---%v\n", srcUrl, "err:AuthorContentUrl is nil", author)
		logrus.Infof("%+v\n\n", author)
	}
}

//通过首页抓取到的作者列表传到这里，这里循环数据去发送请求
func (a *Author) GraspByIndexData(data *define.HomeFormat) {

}

//抓取作者详情信息 /authorv_07d17f8539d7.aspx
func (a *Author) GraspAuthorDetail(authorUrl string, endChan chan bool) {
	defer func() {
		endChan <- true
	}()
	var err error
	logrus.Infoln("GraspAuthorDetail start..........")
	if strings.Contains(authorUrl, "http") == false {
		authorUrl = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(authorUrl, "/")
	}
	if err = a.getSourceHtml(authorUrl, "author.html"); err != nil {
		logrus.Infoln("get url ", authorUrl, "error:", err)
		data.G_GraspResult.PushError(err)
		return
	}
	if a.SourceAuthor.AuthorName == "" {
		a.getAuthorDefaultData()
	}
	a.getAuthorDetailInfo()
	if len(a.SourceAuthor.Data) > 0 {
		sendData := &define.ParseData{
			Data:      a.SourceAuthor,
			Params:    nil,
			ParseFunc: data.NewAuthorStore().LoadAuthorData,
		}
		data.G_GraspResult.SendParseData(sendData)
	}
}

//抓取作者诗词列表数据，并保存诗词列表  /authors/authorvsw_07d17f8539d7A1.aspx
func (a *Author) GraspAuthorPoetryList(authorUrl string, endChan chan bool) {
	defer func() {
		endChan <- true
	}()
	var (
		err        error
		reqUrlList map[uint32]string
	)
	logrus.Infoln("GraspAuthorPoetryList start..........")
	if len(authorUrl) == 0 {
		logrus.Infoln("GraspAuthorPoetryList  err: authorUrl is nil")
		return
	}
	reqUrlList = make(map[uint32]string)
	if strings.Contains(authorUrl, "http") == false {
		authorUrl = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(authorUrl, "/")
	}
	a.Url = authorUrl
	if err = a.getSourceHtml(authorUrl, "authorPoetryList.html"); err != nil {
		logrus.Infoln("GetSourceHtml error:", err)
		return
	}
	logrus.Infoln(authorUrl, "..start....")
	//获取当前页诗词链接信息
	linkMp := a.parsePoetryListLink(a.Html)
	for _, link := range linkMp {
		key := tools.Crc32(link.LinkUrl)
		if _, ok := reqUrlList[key]; ok {
			logrus.Infoln("--", link.LinkUrl, "--已存在")
			continue
		}
		Content.NewContent().GraspContentSaveData(link.LinkUrl, link)
		reqUrlList[key] = link.LinkUrl
		logrus.Infoln("GraspAuthorPoetryList reqUrlList len:", len(reqUrlList))
	}
	a.sendPoetryPageListRequest()
	return
}

//获取诗词列表总页数并发送每页的请求
func (a *Author) sendPoetryPageListRequest() {
	var (
		pageTotalStr string
		totalPageNum int
		err          error
		reqUrlList   map[uint32]string
	)
	reqUrlList = make(map[uint32]string)
	pageTotalStr = a.query.Find(".pagesright>span").Text()
	if totalPageNum, err = tools.TrimAuthorTotalPageText(pageTotalStr); err != nil {
		logrus.Infoln("getPoetryPageList err:", err)
		return
	}
	//最多只能获取10页
	if totalPageNum > 10 {
		totalPageNum = 10
	}
	for i := 2; i <= totalPageNum; i++ {
		iStr := strconv.Itoa(i)
		url := strings.Replace(a.Url, "A1", "A"+iStr, -1)
		go func() {
			logrus.Infoln("url:", url)
			if html, e := base.GetHtml(url); e == nil {
				linkMp := a.parsePoetryListLink(html)
				for _, link := range linkMp {
					key := tools.Crc32(link.LinkUrl)
					if _, ok := reqUrlList[key]; ok {
						logrus.Infoln("--", link.LinkUrl, "--已存在")
						continue
					}
					Content.NewContent().GraspContentSaveData(link.LinkUrl, link)
					reqUrlList[key] = link.LinkUrl
					logrus.Infoln("sendPoetryPageListRequest reqUrlList len :", len(reqUrlList))
				}
			}
		}()
		time.Sleep(50 * time.Millisecond)
	}
}

//解析作者诗词列表页诗词链接
func (a *Author) parsePoetryListLink(html []byte) (linkMap PoetryLinkMap) {
	var (
		query    *goquery.Document
		ok       bool
		linkUrl  string
		linkText string
		err      error
	)
	linkMap = make(PoetryLinkMap)
	if len(html) == 0 {
		return
	}
	if query, err = tools.NewDocumentFromReader(string(html)); err != nil {
		return
	}
	query.Find(".main3>.left>.sons>.cont").Each(func(i int, selection *goquery.Selection) {
		linkUrl, ok = selection.Find("p").Eq(0).Find("a").Attr("href")
		linkText = selection.Find("p").Eq(0).Find("a").Text()
		if ok {
			linkMap[i] = define.LinkStr{
				LinkUrl: linkUrl,
				Title:   linkText,
				Sort:    i,
			}
		}
	})
	return
}

//设置作者信息默认属性值
func (a *Author) SetAuthorAttr(authorSource *define.PoetryAuthorDetail) *Author {
	a.SourceAuthor = authorSource
	return a
}

//获取作者详情页默认数据，如果SetAuthorAttr设置了则不用获取
func (a *Author) getAuthorDefaultData() {
	a.SourceAuthor.AuthorName = strings.TrimSpace(a.query.Find(".main3>.left>.sonspic>.cont>h1>span>b").Text())
	if photoImg, ok := a.query.Find(".main3>.left>.sonspic>.cont>.divimg>img").Attr("src"); ok {
		a.SourceAuthor.AuthorImgUrl = strings.TrimSpace(photoImg)
	}
	totalNumText := a.query.Find(".main3>.left>.sonspic>.cont>p>a").Text()
	if len(totalNumText) > 0 {
		totalNumText = strings.TrimRight(totalNumText, "篇诗文")
		totalNumText = strings.TrimLeft(totalNumText, "►")
		totalNumText = strings.TrimSpace(totalNumText)
		num, _ := strconv.Atoi(totalNumText)
		a.SourceAuthor.AuthorTotalPoetry = num
	}
}

//获取作者详情页作者数据信息
func (a *Author) getAuthorDetailInfo() {
	//作者简介
	if introduction, _ := a.query.Find(".main3>.left>.sonspic>.cont>p").Html(); len(introduction) > 0 {
		index := strings.LastIndex(introduction, "<a")
		a.SourceAuthor.Introduction = introduction[:index] + "</p>"
	}
	//获取资料信息
	a.query.Find(".main3>.left>.sons").Each(func(i int, selection *goquery.Selection) {
		var (
			buf    bytes.Buffer
			detail define.ContentData
		)
		if attrId, ok := selection.Attr("id"); ok && !strings.Contains(attrId, "quan") {
			dataId := strings.TrimLeft(attrId, "fanyi")
			if len(dataId) > 0 {
				detail.DataId, _ = strconv.Atoi(dataId)
				detail.Title = selection.Find(".contyishang>div>h2").Text()
				selection.Find(".contyishang>p").Each(func(i int, selection *goquery.Selection) {
					if html, e := selection.Html(); e == nil {
						buf.WriteString("<p>" + html + "</p>")
					}
				})
				detail.Introd = buf.String()
				buf.Reset()
				dataUrl := config.G_Conf.GuShiWenPoetryUrl + "authors/ajaxziliao.aspx?id=" + dataId
				if bytes, err := base.GetHtml(dataUrl); err == nil {
					detail.Content = string(bytes)
					detail.HtmlSrcUrl = dataUrl
				}
				detail.Sort = i
				detail.Type = int(define.AuthorType)
				detail.PlaySrcUrl = config.G_Conf.GuShiWenPoetryUrl + "ziliaoplay.aspx?id=" + dataId
				detail.PlayUrl = config.G_Conf.GushiwenSongUrl + "machine/ziliao/" + dataId + "/ok.mp3"
				a.SourceAuthor.Data = append(a.SourceAuthor.Data, &detail)
			}
		}
	})
}

//获取页面信息
func (a *Author) getSourceHtml(url string, testFile string) (err error) {
	var (
		bytes []byte
	)
	if config.G_Conf.Env == define.TestEnv {
		//获取测试文件内容
		if testFile == "author.html" {
			return errors.New("stop get file")
		}
		if len(testFile) > 0 {
			dir, _ := os.Getwd()
			file := dir + "/" + testFile
			bytes, err = base.GetTestFile(file)
		} else {
			return errors.New("test file is nil")
		}
	} else {
		bytes, err = base.GetHtml(url)
	}
	if err != nil {
		return
	}
	if len(bytes) > 0 {
		a.Html = bytes
		a.query, err = tools.NewDocumentFromReader(string(bytes))
	}
	return err
}
