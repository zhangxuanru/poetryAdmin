package test

import (
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Author"
	"poetryAdmin/worker/core/grasp/poetry/Category"
	"poetryAdmin/worker/core/grasp/poetry/Content"
	"poetryAdmin/worker/core/parse"
	"runtime"
	"testing"
	"time"
)

//单元测试--抓取全部
func TestGrabsIndex(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	go data.NewGraspResult().PrintMsg()
	go func() {
		for {
			logrus.Infoln("NumGoroutine:", runtime.NumGoroutine())
			time.Sleep(10 * time.Second)
		}
	}()
	var SubReceiveMsg parse.SubscribeMsg
	Convey("测试抓取所有", t, func() {
		SubReceiveMsg = parse.SubscribeMsg{
			PubTile:  "抓取所有",
			AddDate:  time.Now().Unix(),
			Status:   define.TaskStatusImplemented,
			TaskMark: define.GrabPoetryAll,
		}
		parse.NewDispatch(SubReceiveMsg).Execution()
		for {
			time.Sleep(10 * time.Second)
		}
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
		Author.NewAuthor().SendGraspAuthorDataReq(author, "https://so.gushiwen.org/gushi/tangshi.aspx")
	}

	time.Sleep(15 * time.Second)

	//Content.NewContent().GraspContentSaveData("/shiwenv_73add8822103.aspx", nil)
}
