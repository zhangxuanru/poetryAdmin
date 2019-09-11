/*
@Time : 2019/9/9 17:50
@Author : zxr
@File : content
@Software: GoLand
*/
package Action

import (
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/famous/Parser"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"time"
)

//名句获取分类下详情数据
type content struct {
}

func NewContent() *content {
	return &content{}
}

//名句载入主题下分类数据，发送分类请求获取分类下详情数据
//see:https://so.gushiwen.org/mingju/Default.aspx?p=1&c=%e6%8a%92%e6%83%85&t=%e7%88%b1%e6%83%85
func (c *content) LoadThemeCategory(allThemeCategory []*define.ThemeCategory) {
	var (
		themeCategory *define.ThemeCategory
		classify      *define.Classify
		contentList   []define.Content
		page          define.ContentPage
		mp            map[uint32]string
		url           string
		err           error
	)
	mp = make(map[uint32]string)
	//发送分类详情请求，获取具体的数据
	for _, themeCategory = range allThemeCategory {
		for _, classify = range themeCategory.ClassifyList {
			classify.ThemeTitle = themeCategory.Title
			url = classify.LinkUrl

			logrus.Infoln("名句-主题:", themeCategory.Title, "名句-分类:", classify.Title, "start...")
			if _, ok := mp[tools.Crc32(url)]; ok {
				logrus.Infoln("名句分类URL已存在")
				continue
			}
			//todo 循环发送下每页请求
			for page.IsNextPage == true || len(url) > 0 {
				if len(page.NextPageUrl) > 0 {
					url = page.NextPageUrl
				}
				logrus.Infoln("名句-主题:", themeCategory.Title, "名句-分类:", classify.Title, "url:", url)
				if _, ok := mp[tools.Crc32(url)]; ok {
					logrus.Infoln("---名句分类URL已存在")
					continue
				}
				if contentList, page, err = c.callContentPage(url); err != nil {
					logrus.Infoln("callContentNextPage err:", err)
					url = ""
					break
				}
				classify.ContentList = contentList
				//todo 保存contentList数据
				sendData := &define.ParseData{
					Data:      classify,
					Params:    nil,
					ParseFunc: data.NewFamousStorage().LoadClassifyContentData,
				}
				data.G_GraspResult.SendParseData(sendData)
				mp[tools.Crc32(url)] = url
				url = ""
				time.Sleep(50 * time.Millisecond)
			}
			time.Sleep(100 * time.Millisecond)
			mp[tools.Crc32(url)] = url
		}
	}
	return
}

//请求下一页，整理名句列表数据返回
func (c *content) callContentPage(url string) (contentList []define.Content, page define.ContentPage, err error) {
	var bytes []byte
	if bytes, err = c.getUrlSource(url); err != nil {
		return
	}
	if contentList, page, err = Parser.ParseFamousContent(bytes); err != nil {
		logrus.Infoln("ParseFamousContent err:", err)
		return
	}
	return
}

//获取分类详情页数据
func (c *content) getUrlSource(url string) (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		dir, _ := os.Getwd()
		file := dir + "/famous/content1.html"
		return tools.ReadFile(file)
	} else {
		bytes, err = base.GetHtml(url)
	}
	return
}
