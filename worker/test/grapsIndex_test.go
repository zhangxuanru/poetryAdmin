package test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Author"
	"poetryAdmin/worker/core/grasp/poetry/Category"
	"poetryAdmin/worker/core/grasp/poetry/Content"
	"poetryAdmin/worker/core/parse"
	"testing"
	"time"
)

//单元测试--抓取首页
func TestGrabsIndex(t *testing.T) {
	var SubReceiveMsg parse.SubscribeMsg
	Convey("测试抓取所有", t, func() {
		SubReceiveMsg = parse.SubscribeMsg{
			PubTile:  "抓取所有",
			AddDate:  time.Now().Unix(),
			Status:   define.TaskStatusImplemented,
			TaskMark: define.GrabPoetryAll,
		}
		parse.NewDispatch(SubReceiveMsg).Execution()
		time.Sleep(20 * time.Second)
	})
}

//测试诗文类型详情页
func TestCategory(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	home := &define.HomeFormat{
		Identifier: "test",
		Data: define.DataMap{
			1: &define.TextHrefFormat{
				Href:         "https://so.gushiwen.org/gushi/tangshi.aspx",
				Text:         "唐诗三百",
				ShowPosition: 1,
			},
		},
	}
	Category.NewCategory().GraspByIndexData(home)
	time.Sleep(120 * time.Second)
}

//测试诗文详情页
func TestContent(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	poetry := &define.PoetryAuthorList{
		AuthorName:      "柳宗元",
		PoetryTitle:     "江雪",
		PoetrySourceUrl: "/shiwenv_58313be2d918.aspx",
		GenreTitle:      "五言绝句",
		Category: &define.TextHrefFormat{
			Text:         "唐诗三百",
			Href:         "https://so.gushiwen.org/gushi/tangshi.aspx",
			ShowPosition: 1,
		},
	}
	if author := Content.NewContent().GetAuthorContentData(poetry); author.AuthorName != "" {
		Author.NewAuthor().SendGraspAuthorDataReq(author)
	}

	time.Sleep(15 * time.Second)

	//Content.NewContent().GraspContentSaveData("/shiwenv_73add8822103.aspx", nil)
}

func TestA(t *testing.T) {
	format := define.TextHrefFormat{
		Text:         "aa",
		Href:         "aaa",
		ShowPosition: 10,
	}
	logrus.Infof("%+v", format)
	logrus.Infoln(int(format.ShowPosition))
	return

	src := "https://song.gushiwen.org/authorImg/taoyuanming.jpg"
	//src := "https://song.gushiwen.org/machine/ziliao/1601/ok.mp3"
	fileName, err2 := data.NewUploadStore().Upload(src)
	logrus.Infoln(fileName)
	logrus.Infoln(err2)

	return
	file := "D:/server/gitData/goPath/poetryAdmin/worker/test/index.html"
	bytes, err := tools.ReadFile(file)
	logrus.Info("err:", err)
	query, e := tools.NewDocumentFromReader(string(bytes))
	logrus.Info("err:", e)

	query.Find(".right>.sons").Eq(2).Find(".cont>a").Each(func(j int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		logrus.Infoln("href:", href, "text:", selection.Text())
	})
}

func TestB(T *testing.T) {
	ch := make(chan bool, 5)
	end := make(chan bool)
	go func() {
		for {
			select {
			case <-ch:
				logrus.Infoln("ch.......")
			case <-end:
				if len(ch) > 0 {
					logrus.Infoln("还有数据.....")
					continue
				}
				time.Sleep(1 * time.Second)
				goto GoEnd
			}
		}
	GoEnd:
		logrus.Info("end......")
		return
	}()

	for i := 0; i < 10; i++ {
		go func(i int) {
			if i > 4 {
				end <- true
			}
		}(i)
		go func(i int) {
			ch <- true
		}(i)
	}
	time.Sleep(10 * time.Second)
}
