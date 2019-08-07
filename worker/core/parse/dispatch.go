package parse

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/core/data"
	"poetryAdmin/worker/core/define"
	"poetryAdmin/worker/core/engine/poetry"
)

//分发
type Dispatch struct {
	Msg SubscribeMsg
}

func NewDispatch(msg SubscribeMsg) *Dispatch {
	return &Dispatch{
		Msg: msg,
	}
}

//分发执行
func (d *Dispatch) Execution() {
	if d.Msg.TaskMark == "" {
		return
	}
	go data.NewGraspResult().PrintMsg()
	switch d.Msg.TaskMark {
	case define.GrabPoetryAll:
		logrus.Info("Execution start :", d.Msg.PubTile)
		poetry.NewRunAll().Run()
	}
	logrus.Info("end Dispatch Execution.......")
	return
}
