/*
@Time : 2019/8/30 18:14
@Author : zxr
@File : index
@Software: GoLand
*/
package Action

import (
	"errors"
	"os"
	"poetryAdmin/worker/app/config"
	"poetryAdmin/worker/app/tools"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/grasp/ancient/Parser"
	"poetryAdmin/worker/core/grasp/poetry/base"
)

//see https://so.gushiwen.org/guwen/
type Index struct {
}

func NewIndex() *Index {
	return &Index{}
}

func (i *Index) StartGrab() {
	var (
		err   error
		bytes []byte
	)
	if bytes, err = i.getSource(); err != nil {
		data.G_GraspResult.PushError(err)
		return
	}
	Parser.ParseIndex(bytes)
}

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
