package parse

import (
	"github.com/sirupsen/logrus"
	"poetryAdmin/worker/app/config"
)

//接收订阅消息
type Analysis struct {
	ReceiveChan   chan []byte
	SubReceiveMsg SubscribeMsg
}

func NewAnalysis(receiveChan chan []byte) *Analysis {
	return &Analysis{
		ReceiveChan: receiveChan,
	}
}

//接收订阅的频道发送来的消息
func (a *Analysis) ParseSubscribeData() {
	var (
		msg []byte
		err error
	)
	for {
		select {
		case msg = <-a.ReceiveChan:
			if err = a.ReceiveMsgToSubMsg(msg); err != nil {
				continue
			}
			NewDispatch(a.SubReceiveMsg).Execution()
			logrus.Info("sub receive msg:", string(msg))
		}
	}
}

//将获取到的订阅数据json 转为结构体
func (a *Analysis) ReceiveMsgToSubMsg(msg []byte) (err error) {
	subscribeMsg := NewSubscribeMsg()
	err = config.G_Json.Unmarshal(msg, &subscribeMsg)
	logrus.Debugf("subscribeMsg:%+v", subscribeMsg)
	a.SubReceiveMsg = *subscribeMsg
	return err
}
