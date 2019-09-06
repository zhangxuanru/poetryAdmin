/*
@Time : 2019/9/6 18:42
@Author : zxr
@File : index
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
)

//抓取名句首页
//https://so.gushiwen.org/mingju/
type index struct {
}

func NewIndex() *index {
	return &index{}
}

//名句 入口函数
func (i *index) Start() {
	var (
		err          error
		bytes        []byte
		classifyList []define.Classify
	)
	if bytes, err = i.getSource(); err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	if classifyList, err = Parser.ParseFamousIndexClassify(bytes); err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	sendData := &define.ParseData{
		Data:      &classifyList,
		Params:    nil,
		ParseFunc: data.NewFamousStorage().LoadFamousClassifyData,
	}
	data.G_GraspResult.SendParseData(sendData)
	//发送各分类的请求，获取具体内容

}

//获取首页内容
func (i *index) getSource() (bytes []byte, err error) {
	if config.G_Conf.Env == define.TestEnv {
		dir, _ := os.Getwd()
		file := dir + "/famous/index.html"
		return tools.ReadFile(file)
	} else {
		bytes, err = base.GetHtml(config.G_Conf.GushiwenMingJuUrl)
	}
	return
}
