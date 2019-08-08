package data

import (
	"github.com/sirupsen/logrus"
)

//抓取结果处理
type GraspResult struct {
	err   chan error
	close chan bool
	Data  chan interface{}
}

var G_GraspResult *GraspResult

func NewGraspResult() *GraspResult {
	result := &GraspResult{
		err:   make(chan error),
		close: make(chan bool),
		Data:  make(chan interface{}, 5000),
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
func (g *GraspResult) SendData(data interface{}) {
	go func() {
		g.Data <- data
	}()
}

//统一处理错误消息
func (g *GraspResult) PrintMsg() {
	var (
		err   error
		close bool
		data  interface{}
	)
	for {
		select {
		case err = <-g.err:
			logrus.Debug("Execution error:", err)
		case data = <-g.Data:
			logrus.Infoln("ret data:", data)
		case close = <-g.close:
			if close == true {
				goto PRINTERR
			}
		}
	}
PRINTERR:
	logrus.Debug("Execution end:")
	return
}
