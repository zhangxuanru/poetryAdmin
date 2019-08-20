package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/models"
	"poetryAdmin/worker/core/define"
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
