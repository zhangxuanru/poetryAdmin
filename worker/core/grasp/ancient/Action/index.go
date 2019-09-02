/*
@Time : 2019/8/30 18:14
@Author : zxr
@File : index
@Software: GoLand
*/
package Action

import (
	"errors"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/ancient/Parser"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"strings"
	"sync"
)

//see https://so.gushiwen.org/guwen/
type Index struct {
	bookAddr *Book
}

func NewIndex() *Index {
	return &Index{
		bookAddr: NewBook(),
	}
}

//开始抓取古文首页数据
func (i *Index) StartGrab() {
	var (
		err          error
		bytes        []byte
		categoryData []define.GuWenCategoryList
	)
	if bytes, err = i.getSource(); err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	if categoryData, err = Parser.ParseGuWenIndexCategory(bytes); err != nil {
		logrus.Infoln("ParseGuWenIndexCategory err:", err)
		data.G_GraspResult.PushError(err)
		return
	}
	sendData := &define.ParseData{
		Data:      &categoryData,
		Params:    nil,
		ParseFunc: data.NewAncientCategoryStore().LoadCategoryData,
	}
	data.G_GraspResult.SendParseData(sendData)
	go i.bookAddr.ReceiveCategoryBook()
	i.sendCategoryRequest(categoryData)
}

//发送分类请求，获取信书籍息
func (i *Index) sendCategoryRequest(categoryData []define.GuWenCategoryList) {
	var (
		bytes        []byte
		err          error
		cateBookHtml *define.GuWenCategoryBookHtml
		wg           sync.WaitGroup
		categoryNum  int
	)
	for _, category := range categoryData {
		categoryNum += len(category.SubNode)
	}
	wg.Add(categoryNum)
	for _, category := range categoryData {
		for _, cateNode := range category.SubNode {
			if strings.Contains(cateNode.LinkUrl, "http") == false {
				cateNode.LinkUrl = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(cateNode.LinkUrl, "/")
			}
			go func() {
				defer func(wg *sync.WaitGroup) {
					wg.Done()
				}(&wg)
				if bytes, err = i.getCategoryBookSource(cateNode.LinkUrl); err != nil {
					logrus.Infoln("sendCategoryRequest GetHtml err:", err)
					return
				}
				cateBookHtml = &define.GuWenCategoryBookHtml{
					GuWenCategory: define.GuWenCategory{
						CategoryName: cateNode.CategoryName,
						LinkUrl:      cateNode.LinkUrl,
					},
					Html: bytes,
				}
				i.bookAddr.SendCategoryBook(cateBookHtml)
			}()
		}
	}
	wg.Wait()
	i.bookAddr.SendClose(true)
	return
}

//根据分类URL，发送请求
func (i *Index) getCategoryBookSource(url string) (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		dir, _ := os.Getwd()
		file := dir + "/ancient/categoryBook.html"
		return tools.ReadFile(file)
	} else {
		bytes, err = base.GetHtml(url)
	}
	return
}

//古文首页发送HTTP请求
func (i *Index) getSource() (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		bytes, err = i.getTestFile()
	} else {
		bytes, err = base.GetHtml(config.G_Conf.GushiwenAncientUrl)
	}
	return
}

//读取测试的首页文件，避免每次都http请求
func (i *Index) getTestFile() (byt []byte, err error) {
	dir, _ := os.Getwd()
	file := dir + "/ancient/index.html"
	if ret, _ := tools.PathExists(file); ret == true {
		return tools.ReadFile(file)
	}
	return nil, errors.New(file + "file is not exists")
}
