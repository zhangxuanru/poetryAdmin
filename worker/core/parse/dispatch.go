package parse

import (
	"github.com/sirupsen/logrus"
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
	logrus.Info("Execution start :", d.Msg.PubTile)
	switch d.Msg.TaskMark {
	case define.GrabPoetryAll:
		poetry.NewRunAll().Run()
	case define.GrabPoetryRecommend:
		poetry.NewRecommend().Run()
	}
	logrus.Info("end Dispatch Execution.......")
	return
}
