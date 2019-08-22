package test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Category"
	"poetryAdmin/worker/core/grasp/poetry/Content"
	"poetryAdmin/worker/core/parse"
	"strings"
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
		time.Sleep(5 * time.Second)
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
	Content.NewContent().GraspContentData(poetry)

	//Content.NewContent().GraspContentSaveData("/shiwenv_73add8822103.aspx", nil)
}

func TestA(t *testing.T) {
	html := "<p style=\" margin:0px;\">陶渊明（352或365年—427年），字元亮，又名潜，私谥“靖节”，世称靖节先生，浔阳柴桑（今江西省九江市）人。东晋末至南朝宋初期伟大的诗人、辞赋家。曾任江州祭酒、建威参军、镇军参军、彭泽县令等职，最末一次出仕为彭泽县令，八十多天便弃职而去，从此归隐田园。他是中国第一位田园诗人，被称为“古今隐逸诗人之宗”，有《陶渊明集》。<a href=\"/authors/authorvsw_07d17f8539d7A1.aspx\">► 134篇诗文</a></p>"
	index := strings.LastIndex(html, "<a")
	logrus.Infoln(index)
	logrus.Infoln(html[:492])
	logrus.Infoln(html)

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
