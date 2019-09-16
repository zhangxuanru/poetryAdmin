/*
@Time : 2019/9/16 16:13
@Author : zxr
@File : Recommend
@Software: GoLand
*/
package poetry

import (
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Content"
	"poetryAdmin/worker/core/grasp/poetry/Parser"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"strings"
	"time"
)

type Recommend struct {
}

func NewRecommend() *Recommend {
	return &Recommend{}
}

func (r *Recommend) StartGrasp() {
	var (
		bytes       []byte
		contentList []define.PoetryContent
		page        define.ContentPage
		err         error
		count       int
	)
	url := config.G_Conf.GuShiWenIndexUrl
	for (len(page.NextPageUrl) > 0 && count <= 10) || len(url) > 0 {
		if len(page.NextPageUrl) > 0 {
			url = config.G_Conf.GuShiWenIndexUrl + strings.TrimLeft(page.NextPageUrl, "/")
		}
		logrus.Infoln("url:", url)
		if bytes, err = r.GetSource(url); err != nil {
			logrus.Infoln("StartGrasp for GetSource err:", err)
			return
		}
		if contentList, page, err = Parser.ParseRecommendHtml(bytes); err != nil {
			logrus.Infoln("ParseRecommendHtml for err:", err)
			return
		}
		r.ProcContentList(contentList)
		count++
		url = ""
		logrus.Infof("page:%+v\n", page)
	}
}

//处理诗词列表信息
func (r *Recommend) ProcContentList(contentList []define.PoetryContent) {
	var (
		err         error
		contentData models.Content
	)
	for _, content := range contentList {
		if len(content.SourceUrl) == 0 {
			logrus.Infoln("SourceUrl为空......")
			continue
		}
		crc32 := tools.Crc32(content.SourceUrl)
		if contentData, err = models.NewContent().GetContentByCrc32(crc32); err != nil {
			continue
		}
		if contentData.Id == 0 {
			//写入内容表，再写推荐表
			params := define.LinkStr{}
			Content.NewContent().GraspContentSaveData(content.SourceUrl, params)
			contentData, err = models.NewContent().GetContentByCrc32(crc32)
		}
		//写入推荐表
		if contentData.Id > 0 {
			recommend := &models.Recommend{
				PoetryId:    int64(contentData.Id),
				Sort:        content.Sort,
				Status:      1,
				RecommeTime: tools.GetCurrentUnix(),
				AddDate:     time.Now().Unix(),
			}
			if _, err := models.NewRecommend().SaveRecommend(recommend); err != nil {
				logrus.Infoln("SaveRecommend err:", err)
			}
		}
	}
}

func (r *Recommend) GetSource(url string) (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		dir, _ := os.Getwd()
		file := dir + "/recommend.html"
		bytes, err = base.GetTestFile(file)
	} else {
		bytes, err = base.GetHtml(url)
	}
	return
}
