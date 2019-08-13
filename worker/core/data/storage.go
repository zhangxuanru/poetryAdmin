package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/define"
	"time"
)

//保存抓取结果... 分发到各自的module中执行数据库操作
type Storage struct {
	FormatData *define.HomeFormat
	DataMap    define.DataMap
}

func NewStorage() *Storage {
	return new(Storage)
}

//载入data数据
func (s *Storage) LoadData(format *define.HomeFormat) {
	s.FormatData = format
	s.DistributionModule()
}

//载入诗词类型详情页数据
func (s *Storage) LoadCategoryPoetryData(data interface{}, params interface{}) {
	var (
		format    *define.TextHrefFormat
		categorys models.Category
		genreId   int64 //体裁ID
		id        int64 //诗词ID
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
			if categorys, err = models.GetDataByCrcAndCateName(tools.Crc32(format.Href), format.Text, int(format.ShowPosition)); err != nil {
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
				Title:      list.PoetryTitle,
				SourceUrl:  list.PoetrySourceUrl,
				AuthorId:   authorId,
				AddDate:    time.Now().Unix(),
				UpdateDate: time.Now().Unix(),
			}
			if id, err = models.NewContent().SaveContent(content); err != nil {
				G_GraspResult.PushError(err)
				logrus.Debug("SaveContent error:", err)
				continue
			}
			//写入诗词关联表 poetry_content_relation
			relation := &models.ContentRelation{
				PoetryId:   id,
				CategoryId: int64(categorys.Id),
				GenreId:    genreId,
				AuthorId:   authorId,
			}
			if _, err = models.NewContentRelation().SaveContentRelation(relation); err != nil {
				G_GraspResult.PushError(err)
				logrus.Debug("SaveContentRelation error:", err)
			}
		}
	}
	return
}

//分发模块
func (s *Storage) DistributionModule() {
	sign := s.FormatData.Identifier
	var err error
	switch s.FormatData.Identifier {
	case define.HomePoetryCategoryFormatSign: //首页-诗文分类
		_, err = models.InsertMultiCategoryByDataMap(s.formatConversionDataMap())
	case define.HomePoetryFamousFormatSign: //首页-名句
		_, err = models.InsertMultiCategoryByDataMap(s.formatConversionDataMap())
	case define.HomePoetryAuthorFormatSign: //首页-作者
		_, err = models.InsertMultiAuthorByDataMap(s.formatConversionDataMap())
	}
	s.PrintErr(sign, err)
}

//格式转换 dataMap
func (s *Storage) formatConversionDataMap() define.DataMap {
	maps := s.FormatData.Data.(define.DataMap)
	return maps
}

func (s *Storage) PrintErr(sign define.DataFormat, err error) {
	if err != nil {
		logrus.Debugln(sign, "-err:", err)
	}
}
