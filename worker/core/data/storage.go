package data

import (
	"github.com/sirupsen/logrus"
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

//载入数据
func (s *Storage) LoadData(format *define.HomeFormat) {
	s.FormatData = format
}

//分发模块
func (s *Storage) DistributionModule() {
	switch s.FormatData.Identifier {
	case define.HomePoetryCategoryFormatSign:
		maps := s.formatConversionDataMap()

	}
}

//格式转换 dataMap
func (s *Storage) formatConversionDataMap() define.DataMap {
	maps := s.FormatData.Data.(define.DataMap)
	return maps
}
