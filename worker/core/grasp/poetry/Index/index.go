package Index

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"sync"
)

type dataMap map[interface{}]*define.TextHrefFormat

//抓取首页
type Index struct {
	Content      string
	GoQuery      *goquery.Document
	CategoryData dataMap //首页分类数据
	AuthorData   dataMap //首页作者数据
	FamousData   dataMap //首页名句数据
	group        *sync.WaitGroup
}

func NewIndex() *Index {
	return &Index{
		CategoryData: make(dataMap),
		AuthorData:   make(dataMap),
		FamousData:   make(dataMap),
		group:        &sync.WaitGroup{},
	}
}

//获取首页所有内容
func (i *Index) GetAllData() {
	logrus.Info("GetAllData start .......")
	if err := i.GetIndexSource(); err != nil {
		logrus.Debug("GetIndexHtml err:", err)
		return
	}
	if base.CheckContent(i.Content) == false {
		logrus.Debug("CheckContent err: content is nil")
		return
	}
	i.group.Add(3)
	go i.GetPoetryCategory()
	go i.GetPoetryFamousCategory()
	go i.GetPoetryAuthor()
	i.group.Wait()
	return
}

//首页-诗文分类
func (i *Index) GetPoetryCategory() {
	defer i.group.Done()
	if len(i.Content) == 0 || i.GoQuery == nil {
		logrus.Debug("GetPoetryCategory() i.Content is nil or i.query is nil")
		return
	}
	i.GoQuery.Find(".right>.sons:nth-child(1)>.cont>a").Each(func(j int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		result := &define.TextHrefFormat{
			Href: href,
			Text: selection.Text(),
		}
		i.CategoryData[j] = result
	})
	home := &define.HomeFormat{
		Identifier: define.HomePoetryCategoryFormatSign,
		Data:       i.CategoryData,
	}
	data.G_GraspResult.SendData(home)
	return
}

//首页-名句分类
func (i *Index) GetPoetryFamousCategory() {
	defer i.group.Done()
	if len(i.Content) == 0 || i.GoQuery == nil {
		logrus.Debug("GetPoetryFamousCategory() i.Content is nil or i.query is nil")
		return
	}
	i.GoQuery.Find(".right>.sons:nth-child(2)>.cont>a").Each(func(j int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		result := &define.TextHrefFormat{
			Href: href,
			Text: selection.Text(),
		}
		logrus.Infoln("href:", href, "text:", selection.Text())
		i.FamousData[j] = result
	})
	home := &define.HomeFormat{
		Identifier: define.HomePoetryFamousFormatSign,
		Data:       i.FamousData,
	}
	data.G_GraspResult.SendData(home)
	return
}

//首页-作者
func (i *Index) GetPoetryAuthor() {
	defer i.group.Done()
	if len(i.Content) == 0 || i.GoQuery == nil {
		logrus.Debug("GetPoetryAuthor() i.Content is nil or i.query is nil")
		return
	}
	i.GoQuery.Find(".right>.sons:nth-child(3)>.cont>a").Each(func(j int, selection *goquery.Selection) {
		href, _ := selection.Attr("href")
		result := &define.TextHrefFormat{
			Href: href,
			Text: selection.Text(),
		}
		i.AuthorData[j] = result
	})
	home := &define.HomeFormat{
		Identifier: define.HomePoetryAuthorFormatSign,
		Data:       i.AuthorData,
	}
	data.G_GraspResult.SendData(home)
	return
}

//获取首页html内容
func (i *Index) GetIndexSource() (err error) {
	var (
		query *goquery.Document
		bytes []byte
	)
	if config.G_Conf.Env == define.TestEnv {
		bytes, err = i.IndexTestFile()
	} else {
		bytes, err = base.GetHtml(config.G_Conf.GuShiWenIndexUrl)
	}
	if err != nil {
		return
	}
	if len(bytes) > 0 {
		i.Content = string(bytes)
		query, err = tools.NewDocumentFromReader(i.Content)
	}
	if err != nil {
		return err
	}
	i.GoQuery = query
	return nil
}

//读取测试的首页文件，避免每次都http请求
func (i *Index) IndexTestFile() (byt []byte, err error) {
	dir, _ := os.Getwd()
	file := dir + "/index.html"
	if ret, _ := tools.PathExists(file); ret == true {
		return tools.ReadFile(file)
	}
	return nil, errors.New(file + "file is not exists")
}
