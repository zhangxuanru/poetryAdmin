package Category

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/poetry/Author"
	"poetryAdmin/worker/core/grasp/poetry/Content"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"strings"
	"sync"
	"time"
)

//诗文分类模块 抓取诗文分类
type Category struct {
	wg  *sync.WaitGroup
	url string
}

func NewCategory() *Category {
	return &Category{
		wg: &sync.WaitGroup{},
	}
}

//通过首页抓取到的诗文分类传到这里，这里循环数据去发送请求
func (c *Category) GraspByIndexData(data *define.HomeFormat) {
	datas := data.Data.(define.DataMap)
	if len(datas) == 0 {
		return
	}
	c.wg.Add(len(datas))
	for _, ret := range datas {
		c.GetCategorySource(ret.Href, ret)
		time.Sleep(2 * time.Millisecond)
	}
	c.wg.Wait()
}

//获取诗文分类体裁和作者过虑
func (c *Category) GetCategorySource(url string, category *define.TextHrefFormat) {
	defer c.wg.Done()
	var (
		bytes []byte
		err   error
	)
	if config.G_Conf.Env == define.TestEnv {
		bytes, err = c.CateTestFile()
	} else {
		bytes, err = base.GetHtml(url)
	}
	if err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	logrus.Infoln("GetCategorySource url:", url)
	c.url = url
	c.FindDocument(bytes, category)
	return
}

//诗文详情页页面内容分析
func (c *Category) FindDocument(bytes []byte, category *define.TextHrefFormat) (dataMap define.PoetryDataMap) {
	var (
		query *goquery.Document
		err   error
	)
	query, err = tools.NewDocumentFromReader(string(bytes))
	if err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	dataMap = make(define.PoetryDataMap)
	query.Find(".left>.sons").Eq(0).Find(".typecont").Each(func(i int, selection *goquery.Selection) {
		genreTitle := selection.Find(".bookMl").Text()
		if genreTitle == "" || len(genreTitle) == 0 {
			genreTitle = "无"
		}
		selection.Find("span").Each(func(j int, selection *goquery.Selection) {
			poetryText := selection.Text()
			href, _ := selection.Find("a").Attr("href")
			if strings.Contains(href, "http") == false {
				href = config.G_Conf.GuShiWenPoetryUrl + strings.TrimLeft(href, "/")
			}
			if len(poetryText) > 0 {
				splitArr := strings.Split(poetryText, "(")
				AuthorName := ""
				if len(splitArr) > 1 {
					AuthorName = strings.TrimRight(splitArr[1], ")")
				}
				poetryAuthors := &define.PoetryAuthorList{
					AuthorName:        AuthorName,
					PoetryTitle:       splitArr[0],
					PoetrySourceUrl:   href,
					CategoryAuthorUrl: c.url,
					GenreTitle:        genreTitle,
					Category:          category,
				}
				dataMap[genreTitle] = append(dataMap[genreTitle], poetryAuthors)
			}
		})
	})
	sendData := &define.ParseData{
		Data:      &dataMap,
		Params:    category,
		ParseFunc: data.NewCategoryStorage().LoadCategoryPoetryData,
	}
	data.G_GraspResult.SendParseData(sendData)
	c.goPoetryDetail(&dataMap)
	return dataMap
}

//发送进入诗文详情页的请求
func (c *Category) goPoetryDetail(dataMap *define.PoetryDataMap) {
	//对作者进行过虑，同一个作者只发一次请求到详情页，
	//如果没有作者， 则也进入详情页
	var (
		sysMap sync.Map
	)
	for _, ret := range *dataMap {
		for k, val := range ret {
			list := val.(*define.PoetryAuthorList)
			key := fmt.Sprintf("author%d", k)
			if len(list.AuthorName) > 0 {
				sysMap.Store(list.AuthorName, list)
			} else {
				sysMap.Store(key, list)
			}
		}
	}
	//过虑后发送请求
	authorListMp := make(map[string]*define.PoetryAuthorList)
	authorList := []*define.PoetryAuthorList{}
	sysMap.Range(func(key, value interface{}) bool {
		if val, ok := value.(*define.PoetryAuthorList); ok {
			if val.AuthorName != "" {
				authorListMp[val.AuthorName] = val
			} else {
				authorList = append(authorList, val)
			}
		}
		return true
	})
	for _, poetryAuthor := range authorListMp {
		if author := Content.NewContent().GetAuthorContentData(poetryAuthor); author != nil && author.AuthorName != "" {
			logrus.Infoln("开始抓取：", author.AuthorName, "的诗词,来源页:", c.url, "--start")
			Author.NewAuthor().SendGraspAuthorDataReq(author, c.url)
			logrus.Infoln("抓取：", author.AuthorName, "的诗词,来源页:", c.url, "--end")
		}
	}
	for _, poetryAuthor := range authorList {
		if author := Content.NewContent().GetAuthorContentData(poetryAuthor); author != nil && author.AuthorName != "" {
			logrus.Infoln("开始抓取：", author.AuthorName, "的诗词,来源页:", c.url, "--start")
			Author.NewAuthor().SendGraspAuthorDataReq(author, c.url)
			logrus.Infoln("抓取：", author.AuthorName, "的诗词,来源页:", c.url, "--end")
		}
	}
	//sysMap.Range(func(key, value interface{}) bool {
	//	val := value.(*define.PoetryAuthorList)
	//	go func() {
	//		if author := Content.NewContent().GetAuthorContentData(val); author.AuthorName != "" {
	//			Author.NewAuthor().SendGraspAuthorDataReq(author)
	//		}
	//	}()
	//	randI := time.Duration(tools.RandInt64(50, 300))
	//	time.Sleep(randI * time.Millisecond)
	//	return true
	//})
}

//读取测试文件内容
func (c *Category) CateTestFile() (byt []byte, err error) {
	dir, _ := os.Getwd()
	file := dir + "/category.html"
	if ret, _ := tools.PathExists(file); ret == true {
		return tools.ReadFile(file)
	}
	return nil, errors.New(file + "file is not exists")
}
