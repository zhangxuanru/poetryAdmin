/*
@Time : 2019/9/11 11:23
@Author : zxr
@File : store_test
@Software: GoLand
*/
package test

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/ancient/Entrance"
	"poetryAdmin/worker/core/grasp/famous/Action"
	"poetryAdmin/worker/core/grasp/poetry/Author"
	"poetryAdmin/worker/core/grasp/poetry/Content"
	"poetryAdmin/worker/core/grasp/poetry/Parser"
	"poetryAdmin/worker/core/grasp/poetry/base"
	"testing"
)

//测试保存名句
func TestStoreFamous(t *testing.T) {
	var (
		err              error
		categoryMap      map[uint32]models.Category
		mps              map[int][]models.Category
		cateList         []models.Category
		ok               bool
		allThemeCategory []*define.ThemeCategory
	)
	if categoryMap, err = models.GetCategoryDataByPosition(2); err != nil {
		logrus.Infoln("GetCategoryDataByPosition err:", err)
	}
	mps = make(map[int][]models.Category)
	for _, category := range categoryMap {
		if category.Pid > 0 {
			mps[category.Pid] = append(mps[category.Pid], category)
		}
	}
	for _, category := range categoryMap {
		var (
			classifyList  []*define.Classify
			themeCategory define.ThemeCategory
		)
		if category.Pid > 0 {
			continue
		}
		themeCategory.Title = category.CatName
		themeCategory.LinkUrl = category.SourceUrl
		if cateList, ok = mps[category.Id]; ok == false {
			continue
		}
		for _, cateVal := range cateList {
			classify := &define.Classify{
				Title:   cateVal.CatName,
				LinkUrl: cateVal.SourceUrl,
			}
			classifyList = append(classifyList, classify)
		}
		themeCategory.ClassifyList = classifyList
		allThemeCategory = append(allThemeCategory, &themeCategory)
	}
	Action.NewContent().LoadThemeCategory(allThemeCategory)
}

//测试保存古籍
func TestStoreAncient(t *testing.T) {
	Entrance.NewGrab().Exec()
}

//作者source_url为空的情况，然后保存诗词
func TestStorePoetry(t *testing.T) {
	var (
		authors      []models.Author
		authorDetail define.PoetryAuthorDetail
		err          error
		bytes        []byte
	)
	go data.NewGraspResult().PrintMsg()
	//把没有source_url的作者查出来，去请求搜索然后补到 表里
	_, err = models.NewAuthor().GetOrm().QueryTable(models.TableAuthor).Filter("source_url", "").All(&authors, "author", "id")
	if err != nil {
		logrus.Infoln("setup 1: err:", err)
	}
	for _, author := range authors {
		searchUrl := fmt.Sprintf("https://so.gushiwen.org/search.aspx?value=%s", author.Author)
		if bytes, err = getUrlHtml(searchUrl); err != nil {
			logrus.Infoln("getUrlHtml err:", err)
			continue
		}
		if authorDetail, err = Parser.ParserSearchAuthorHtml(bytes); err != nil {
			logrus.Infoln("ParserSearchAuthorHtml err:", err)
			continue
		}
		if len(authorDetail.AuthorSrcUrl) > 0 {
			Author.NewAuthor().SetAuthorAttr(&authorDetail).GraspAuthorDetail(authorDetail.AuthorSrcUrl)
		}
		if len(authorDetail.AuthorContentUrl) > 0 {
			Author.NewAuthor().SetAuthorAttr(&authorDetail).GraspAuthorPoetryList(authorDetail.AuthorContentUrl)
		}
	}
}

//poetry_content表content为空的情况 保存信息
func TestStoreAuthorPoetry(t *testing.T) {
	var (
		err    error
		max    int
		start  int
		count  int64
		count1 int64
	)
	start = 25000
	max = 40000
	maxId := 0
	for start < max {
		logrus.Infoln("start=", start, "-max=", max)
		var dataAll []models.Content
		_, err = models.NewContent().GetOrm().QueryTable(models.TableContent).Filter("id__gte", start).Limit(20).OrderBy("id").All(&dataAll, "id", "content", "source_url")
		if err != nil {
			logrus.Infoln("err:", err)
			return
		}
		for _, content := range dataAll {
			logrus.Infoln("SourceUrl:", content.SourceUrl)
			if content.Id > maxId {
				maxId = content.Id
			}
			//如果诗词表内容为空，则补录诗词信息
			if len(content.Content) == 0 {
				params := define.LinkStr{}
				Content.NewContent().GraspContentSaveData(content.SourceUrl, params)
			} else {
				if count, err = models.NewContentRec().GetOrm().QueryTable(models.TableRec).Filter("poetry_id", content.Id).Count(); err != nil {
					logrus.Infoln("NewContentRec get error:", err)
					continue
				}
				if count1, err = models.NewContentTrans().GetOrm().QueryTable(models.TableTrans).Filter("poetry_id", content.Id).Count(); err != nil {
					logrus.Infoln("NewContentTrans get error:", err)
					continue
				}
				if count+count1 == 0 {
					params := define.LinkStr{}
					Content.NewContent().GraspContentSaveData(content.SourceUrl, params)
				}
			}
			logrus.Infoln("处理:", content.Id, "结束......")
		}
		start = maxId
	}
}

func getUrlHtml(url string) (bytes []byte, err error) {
	bytes, err = base.GetHtml(url)
	return
}
