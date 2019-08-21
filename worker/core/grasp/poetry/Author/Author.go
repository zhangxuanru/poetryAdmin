package Author

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"strconv"
	"strings"
)

//作者模块  抓取作者
type Author struct {
	SourceAuthor *define.PoetryAuthorDetail
	Html         []byte
	query        *goquery.Document
}

func NewAuthor() *Author {
	return &Author{}
}

//通过首页抓取到的作者列表传到这里，这里循环数据去发送请求
func (a *Author) GraspByIndexData(data *define.HomeFormat) {

}

//获取作者详情信息 /authorv_07d17f8539d7.aspx
func (a *Author) GetAuthorDetail(authorUrl string, endChan chan bool) {
	defer func() {
		endChan <- true
	}()
	var (
		err error
	)
	if strings.Contains(authorUrl, "http:") == false {
		authorUrl = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(authorUrl, "/")
	}
	if err = a.getSourceHtml(authorUrl); err != nil {
		logrus.Infoln("get url ", authorUrl, "error:", err)
		data.G_GraspResult.PushError(err)
		return
	}
	if a.SourceAuthor.AuthorName == "" {
		a.getAuthorDefaultData()
	}
	a.getAuthorDetailInfo()
}

//获取作者诗词列表数据，并保存诗词列表
func (a *Author) GetAuthorPoetryList(authorUrl string, endChan chan bool) {
	defer func() {
		endChan <- true
	}()
	logrus.Infoln("authorUrl:", authorUrl)

	//var (
	//	err error
	//)
	//url := ""
	//if err = a.getSourceHtml(url); err != nil {
	//	logrus.Infoln("GetSourceHtml error:", err)
	//	return
	//}
	//a.GetAuthorDefaultData()
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
	a.SourceAuthor.Introduction, _ = a.query.Find(".main3>.left>.sonspic>.cont>p").Html()
	//-------未完待续-------
}

//获取页面信息
func (a *Author) getSourceHtml(url string) (err error) {
	var (
		bytes []byte
	)
	if config.G_Conf.Env == define.TestEnv {
		//获取测试文件内容
		dir, _ := os.Getwd()
		file := dir + "/author.html"
		bytes, err = base.GetTestFile(file)
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
