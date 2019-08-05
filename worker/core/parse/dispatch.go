package parse

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/core/define"
)

//分发
type Dispatch struct {
	Msg SubscribeMsg
}

func NewDispatch(msg SubscribeMsg) *Dispatch {
	return &Dispatch{Msg: msg}
}

//分发执行
func (d *Dispatch) Execution() {
	if d.Msg.TaskMark == "" {
		return
	}
	switch d.Msg.TaskMark {
	case define.GrabPoetryAll:
		logrus.Info("Execution:", d.Msg.PubTile)
		logrus.Info("Execution 开始执行, 执行所有抓取")

	}
}
