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
