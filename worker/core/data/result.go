package data

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/core/define"
)

const ChanMaxLen = 5000

//抓取结果处理
type GraspResult struct {
	err   chan error
	close chan bool
	Data  chan *define.HomeFormat
}

var G_GraspResult *GraspResult

func NewGraspResult() *GraspResult {
	result := &GraspResult{
		err:   make(chan error),
		close: make(chan bool),
		Data:  make(chan *define.HomeFormat, ChanMaxLen),
	}
	G_GraspResult = result
	return result
}

//发送错误消息
func (g *GraspResult) PushError(err error) {
	g.err <- err
}

//发送错误消息并关闭协和
func (g *GraspResult) PushErrorAndClose(err error) {
	g.PushError(err)
	g.PushCloseMark(true)
}

//发送结束标志信息
func (g *GraspResult) PushCloseMark(mark bool) {
	g.close <- mark
}

//发送数据
func (g *GraspResult) SendData(data *define.HomeFormat) {
	g.Data <- data
}

//统一处理错误消息
func (g *GraspResult) PrintMsg() {
	var (
		err       error
		close     bool
		data      interface{}
		autoClose bool
	)
	for {
		if autoClose == true && len(g.Data) == 0 {
			goto PRINTERR
		}
		select {
		case err = <-g.err:
			logrus.Debug("Execution error:", err)
		case data = <-g.Data:
			logrus.Infoln("ret data:", data)
		case close = <-g.close:
			if len(g.Data) > 0 {
				autoClose = true
				logrus.Info("data 还有数据，暂时不能退出")
				continue
			}
			if close == true {
				goto PRINTERR
			}
		}
	}
PRINTERR:
	logrus.Debug("Execution end:")
	return
}
