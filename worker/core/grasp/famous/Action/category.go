/*
@Time : 2019/9/9 14:38
@Author : zxr
@File : content
@Software: GoLand
*/
package Action

import (
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/famous/Parser"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"strings"
)

//名句- 获取分类详情数据
//see:https://so.gushiwen.org/mingju/Default.aspx?p=1&c=%e6%8a%92%e6%83%85
type category struct {
}

func NewCategory() *category {
	return &category{}
}

//载入主题数据，获取分类内容
func (c *category) LoadClassifyList(list *[]define.Classify) {
	var (
		bytes            []byte
		classifyList     []*define.Classify
		allThemeCategory []*define.ThemeCategory
		themeCategory    define.ThemeCategory
		err              error
	)
	for _, classify := range *list {
		url := classify.LinkUrl
		if strings.Contains(url, "http") == false {
			url = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(url, "/")
		}
		if bytes, err = c.getUrlSource(url); err != nil {
			data.G_GraspResult.PushError(err)
			continue
		}
		if classifyList, err = Parser.ParseFamousThemeClassify(bytes); err != nil {
			data.G_GraspResult.PushError(err)
			continue
		}
		themeCategory.Title = classify.Title
		themeCategory.LinkUrl = url
		themeCategory.ClassifyList = classifyList
		sendData := &define.ParseData{
			Data:      &themeCategory,
			Params:    nil,
			ParseFunc: data.NewFamousStorage().LoadFamousCategoryData,
		}
		data.G_GraspResult.SendParseData(sendData)
		allThemeCategory = append(allThemeCategory, &themeCategory)
	}
	NewContent().LoadThemeCategory(allThemeCategory)
}

//获取主题详情页数据
func (c *category) getUrlSource(url string) (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		dir, _ := os.Getwd()
		file := dir + "/famous/content.html"
		return tools.ReadFile(file)
	} else {
		bytes, err = base.GetHtml(url)
	}
	return
}
