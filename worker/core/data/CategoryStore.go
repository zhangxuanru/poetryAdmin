package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"time"
)

//保存分类结果...
type categoryStorage struct {
}

func NewCategoryStorage() *categoryStorage {
	return new(categoryStorage)
}

//载入诗词类型详情页数据
func (c *categoryStorage) LoadCategoryPoetryData(data interface{}, params interface{}) {
	var (
		format    *define.TextHrefFormat
		categorys models.Category
		genreId   int64 //体裁ID
		authorId  int64
		err       error
		ok        bool
	)
	dataMap := data.(*define.PoetryDataMap)
	if len(*dataMap) == 0 {
		return
	}
	format, ok = params.(*define.TextHrefFormat)
	for genreTitle, authorList := range *dataMap {
		if genreTitle != "无" && ok {
			if categorys, err = models.GetDataByCateName(format.Text, int(format.ShowPosition)); err != nil {
				G_GraspResult.PushError(err)
				continue
			}
			//保存 诗文体裁
			gen := &models.Genre{
				GenreName: genreTitle.(string),
				AddDate:   time.Now().Unix(),
			}
			if genreId, _ = models.SaveGenre(gen); genreId > 0 {
				//保存诗文类别体裁关联表
				cateGem := &models.CategoryGenre{
					CatId:   categorys.Id,
					GenreId: genreId,
				}
				_, err = models.NewCategoryGenre().SaveCategoryGenre(cateGem)
			}
		}
		for _, author := range authorList {
			list := author.(*define.PoetryAuthorList)
			//写入作者表
			author := &models.Author{
				Author: list.AuthorName,
			}
			if authorId, err = models.NewAuthor().SaveAuthor(author); err != nil {
				G_GraspResult.PushError(err)
				logrus.Debug("SaveAuthor error:", err)
			}
			//写入诗词表 poetry_content
			content := &models.Content{
				Title:          list.PoetryTitle,
				SourceUrl:      list.PoetrySourceUrl,
				SourceUrlCrc32: tools.Crc32(list.PoetrySourceUrl),
				AuthorId:       authorId,
				GenreId:        genreId,
				AddDate:        time.Now().Unix(),
				UpdateDate:     time.Now().Unix(),
			}
			if _, err = models.NewContent().SaveContent(content); err != nil {
				go G_GraspResult.PushError(err)
				logrus.Debug("SaveContent error:", err)
				continue
			}
		}
	}
	logrus.Infoln("LoadCategoryPoetryData ok.......")
	return
}
