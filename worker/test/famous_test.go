/*
@Time : 2019/9/6 19:07
@Author : zxr
@File : famous_test
@Software: GoLand
*/
package test

import (
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/famous/Action"
	"poetryAdmin/worker/core/grasp/famous/Parser"
	"strings"
	"testing"
	"time"
)

//抓取名句首页
func TestFamousIndex(t *testing.T) {
	go data.NewGraspResult().PrintMsg()

	Action.NewIndex().Start()

	time.Sleep(5 * time.Second)
}

//根据主题抓取名句分类
func TestFamousContentByClassify(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	classify := &[]define.Classify{
		{
			Title:   "抒情",
			LinkUrl: "/mingju/Default.aspx?p=1&c=%e6%8a%92%e6%83%85",
		},
		//{
		//	Title:   "山水",
		//	LinkUrl: "/mingju/Default.aspx?p=1&c=%e5%b1%b1%e6%b0%b4",
		//},
	}
	Action.NewCategory().LoadClassifyList(classify)
	time.Sleep(5 * time.Second)
}

//测试抓取具体名句内容
func TestFamousContent(t *testing.T) {
	go data.NewGraspResult().PrintMsg()
	allThemeCategory := []*define.ThemeCategory{
		&define.ThemeCategory{
			Title:   "抒情",
			LinkUrl: "https://so.gushiwen.org/mingju/Default.aspx?p=1&c=%e6%8a%92%e6%83%85",
			ClassifyList: []*define.Classify{
				&define.Classify{
					ThemeTitle: "抒情",
					Title:      "思乡",
					LinkUrl:    "https://so.gushiwen.org/mingju/Default.aspx?p=1&c=%e6%8a%92%e6%83%85&t=%e6%80%9d%e4%b9%a1",
				},
			},
		},
	}
	Action.NewContent().LoadThemeCategory(allThemeCategory)
	time.Sleep(5 * time.Second)
}

func TestStr(t *testing.T) {
	str := "严羽《满江红·送廖叔仁赴阙》"
	split := strings.Split(str, "《")
	logrus.Infof("%s---%s", split[0], split[1])

	index := strings.Index(str, "《")
	logrus.Infoln("index:", index)

	logrus.Infoln(str[:6])

	logrus.Infoln(str[6:])
}

func TestAA(t *testing.T) {
	var (
		bytes []byte
		err   error
	)
	dir, _ := os.Getwd()
	file := dir + "/famous/content1.html"
	if bytes, err = tools.ReadFile(file); err != nil {
		logrus.Infoln("err:", err)
		return
	}
	contentList, page, err := Parser.ParseFamousContent(bytes)
	logrus.Infoln("contentList:", contentList)
	logrus.Infoln("page:", page)
	logrus.Infoln("err:", err)

}
